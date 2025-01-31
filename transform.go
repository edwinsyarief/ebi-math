package ebimath

type Transformer interface {
	GetParentTransform() *Transform
	GetTransform() *Transform
}

type Transform struct {
	position, scale, offset, origin      Vector
	rotation                             float64
	parent                               *Transform
	dirty, parentDirty                   bool
	matrix, parentMatrix, parentInverted Matrix
}

func (self *Transform) GetParentTransform() *Transform {
	return self.parent
}

// GetInitialParentTransform finds the topmost parent in the hierarchy.
func (self *Transform) GetInitialParentTransform() *Transform {
	for self.parent != nil {
		self = self.parent
	}
	return self
}

func (self *Transform) GetTransform() *Transform {
	return self
}

// T creates a new Transform with default values.
func T() *Transform {
	return &Transform{
		scale: V2(1),
	}
}

func (self *Transform) Origin() Vector {
	return self.origin
}

func (self *Transform) SetOrigin(origin Vector) {
	self.dirty = true
	self.parentDirty = true
	self.origin = origin
}

func (self *Transform) IsDirty() bool {
	return self.dirty || self.parentDirty
}

// SetPosition updates the position, considering parent transforms.
func (self *Transform) SetPosition(position Vector) {
	self.dirty = true
	self.parentDirty = true
	if self.parent != nil {
		_, mi := self.parent.MatrixForParenting()
		self.position = position.Apply(mi)
	} else {
		self.position = position
	}
}

func (self *Transform) SetRotation(rotation float64) {
	self.dirty = true
	self.parentDirty = true
	if self.parent != nil {
		self.rotation -= self.parent.Rotation()
	} else {
		self.rotation = rotation
	}
}

func (self *Transform) SetScale(scale Vector) {
	self.dirty = true
	self.parentDirty = true
	self.scale = scale
}

func (self *Transform) AddScale(add ...Vector) {
	self.dirty = true
	self.parentDirty = true
	self.scale = self.scale.Add(add...)
}

func (self *Transform) SetOffset(offset Vector) {
	self.dirty = true
	self.parentDirty = true
	self.offset = offset
}

// Abs returns an absolute transform without considering parents.
func (self *Transform) Abs() Transform {
	if self.parent == nil {
		return *self
	}
	abs := *T()
	abs.position = self.Position()
	abs.rotation = self.Rotation()
	abs.origin = self.Origin()
	abs.scale = self.Scale()
	abs.offset = self.Offset()
	abs.dirty = true
	abs.parentDirty = true
	return abs
}

func (self *Transform) Rel() Transform {
	rel := *self
	rel.parent = nil
	return rel
}

func (self *Transform) Position() Vector {
	if self.parent == nil {
		return self.position
	}
	pm, _ := self.parent.MatrixForParenting()
	return self.position.Apply(pm)
}

func (self *Transform) Move(v ...Vector) {
	self.SetPosition(self.Position().Add(v...))
}

func (self *Transform) Scale() Vector {
	return self.scale
}

func (self *Transform) Rotation() float64 {
	if self.parent == nil {
		return self.rotation
	}
	return self.rotation + self.parent.Rotation()
}

func (self *Transform) Rotate(rotation float64) {
	self.dirty = true
	self.parentDirty = true
	self.rotation += rotation
}

func (self *Transform) Offset() Vector {
	return self.offset
}

func (self *Transform) Connected() bool {
	return self.parent != nil
}

func (self *Transform) Replace(new Transformer) {
	nt := new.GetTransform()
	self.SetPosition(nt.Position())
	self.SetOffset(nt.Offset())
	self.SetRotation(nt.Rotation())
	self.SetOrigin(nt.Origin())
}

// Connect establishes a parent-child relationship.
func (self *Transform) Connect(parent Transformer) {
	if parent == nil {
		return
	}
	self.parent = parent.GetTransform()
	self.SetPosition(self.Position())
	self.SetRotation(self.Rotation())
	self.SetScale(self.Scale())
	self.SetOffset(self.Offset())
}

// Disconnect removes the parent relationship, making the transform absolute.
func (self *Transform) Disconnect() {
	if self.parent == nil {
		return
	}
	*self = self.Abs()
}

// MatrixForParenting returns matrices for child positioning.
func (self *Transform) MatrixForParenting() (Matrix, Matrix) {
	if self.parentDirty {
		self.parentMatrix = Matrix{}
		self.parentMatrix.Scale(self.scale.X, self.scale.Y)
		self.parentMatrix.Translate(-self.offset.X*self.scale.X, -self.offset.Y*self.scale.Y)
		self.parentMatrix.Rotate(self.rotation)
		self.parentMatrix.Translate(self.position.X, self.position.Y)

		// Store inverted matrix for efficiency
		self.parentInverted = self.parentMatrix
		self.parentInverted.Invert()
		self.parentDirty = false
	}

	m, mi := self.parentMatrix, self.parentInverted
	if self.parent != nil {
		pm, pmi := self.parent.MatrixForParenting()
		m.Concat(pm)
		mi = pmi
		mi.Concat(self.parentInverted)
	}
	return m, mi
}

// Matrix computes the transformation matrix for this node.
func (self *Transform) Matrix() Matrix {
	if self.dirty {
		self.matrix = Matrix{}
		self.matrix.Scale(self.scale.X, self.scale.Y)
		self.matrix.Translate(-self.offset.X*self.scale.X, -self.offset.Y*self.scale.Y)
		self.matrix.Rotate(self.rotation)
		self.matrix.Translate(self.position.X, self.position.Y)
		if !self.origin.IsZero() {
			self.matrix.Translate(-self.origin.X, -self.origin.Y)
		}
		self.dirty = false
	}

	m := self.matrix
	if self.parent != nil {
		pm, _ := self.parent.MatrixForParenting()
		m.Concat(pm)
	}
	return m
}
