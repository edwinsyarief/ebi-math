package ebimath

type Transformer interface {
	GetParentTransform() *Transform
	GetTransform() *Transform
}

// The structure represents basic transformation
// features: positioning, rotating and scaling.
type Transform struct {
	// Absolute (if no parent) position and
	// the scale.
	position, scale Vector
	// The object rotation in radians.
	rotation float64
	// The not scaled offset vector from upper left corner
	// which the object should be rotated around.
	offset Vector
	// origin offset (camera relative position)
	origin Vector

	// If is not nil then the upper values will be relational to
	// the parent ones.
	parent *Transform

	// Dirty is true if we anyhow changed matrix.
	dirty, parentDirty bool

	matrix, parentMatrix, parentInverted Matrix
}

func (self *Transform) GetParentTransform() *Transform {
	return self.parent
}

func (self *Transform) GetInitialParentTransform() *Transform {
	parent := self.GetParentTransform()

	for parent != nil {
		nextParent := parent.GetParentTransform()
		if nextParent == nil {
			break
		}
		parent = nextParent
	}

	return parent
}

// For implementing the Transformer on embedding.
func (self *Transform) GetTransform() *Transform {
	return self
}

// Returns the default Transform structure.
func T() *Transform {
	ret := &Transform{
		position: V2(0),
		scale:    V2(1),
		offset:   V2(0.0),
		rotation: 0.0,
	}
	return ret
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

// Set the absolute object position.
func (self *Transform) SetPosition(position Vector) {
	self.dirty = true
	self.parentDirty = true
	if self.parent != nil {
		_, mi := self.parent.MatrixForParenting()
		self.position = position.Apply(mi)
		return
	}
	self.position = position
}

// Set the absolute object rotation.
func (self *Transform) SetRotation(rotation float64) {
	self.dirty = true
	self.parentDirty = true
	if self.parent != nil {
		self.rotation -= self.parent.Rotation()
		return
	}
	self.rotation = rotation
}

// Set the absolute object scale.
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

// Get the absolute representation of the transform.
func (self *Transform) Abs() Transform {
	if self.parent == nil {
		return *self
	}

	ret := *T()
	ret.position = self.Position()
	ret.rotation = self.Rotation()
	ret.origin = self.Origin()
	ret.scale = self.Scale()
	ret.offset = self.Offset()
	ret.dirty = true
	ret.parentDirty = true

	return ret
}

func (self *Transform) Rel() Transform {
	ret := *self
	ret.parent = nil
	return ret
}

// Get the absolute object position.
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

// Get the absolute object scale.
func (self *Transform) Scale() Vector {
	return self.scale
}

// Get the absolute object rotation.
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

// Returns true if the object is connected
// to some parent.
func (self *Transform) Connected() bool {
	return self.parent != nil
}

func (self *Transform) Replace(new Transformer) {
	self.SetPosition(new.GetTransform().Position())
	self.SetOffset(new.GetTransform().Offset())
	self.SetRotation(new.GetTransform().Rotation())
	self.SetOrigin(new.GetTransform().Origin())
}

// Connect the object to another one making it its parent.
func (self *Transform) Connect(parent Transformer) {
	if parent == nil {
		return
	}
	if self.parent != nil {
		self.Disconnect()
	}

	self.parent = parent.GetTransform()
	self.SetPosition(self.Position()) // Update position based on new parent
	self.SetRotation(self.Rotation()) // Update rotation based on new parent
	self.SetScale(self.Scale())       // Maintain scale
	self.SetOffset(self.Offset())     // Maintain offset
}

// Disconnect from the parent.
func (self *Transform) Disconnect() {
	if self.parent == nil {
		return
	}
	*self = self.Abs()
}

// Return the matrix and the inverted one for parenting children.
func (self *Transform) MatrixForParenting() (Matrix, Matrix) {
	var m, mi Matrix
	if self.parentDirty {
		// Scale first.
		m.Scale(self.scale.X, self.scale.Y)
		// Then move and rotate.
		m.Translate(
			-self.offset.X*self.scale.X,
			-self.offset.Y*self.scale.Y,
		)
		m.Rotate(float64(self.rotation))
		m.Translate(self.position.X, self.position.Y)
		self.parentMatrix = m

		mi = m
		mi.Invert()
		self.parentInverted = mi

		self.parentDirty = false
	} else {
		m = self.parentMatrix
		mi = self.parentInverted
	}

	if self.parent != nil {
		pm, pmi := self.parent.MatrixForParenting()
		m.Concat(pm)
		pmi.Concat(mi)
		mi = pmi
	}

	return m, mi

}

// Returns the GeoM with corresponding
// to the transfrom transformation.
func (self *Transform) Matrix() Matrix {
	var m, pm Matrix

	// Calculating only if we changed the structure anyhow.
	if self.dirty {
		// Scale first.
		m.Scale(self.scale.X, self.scale.Y)
		// Then move and rotate.
		m.Translate(
			-self.offset.X*self.scale.X,
			-self.offset.Y*self.scale.Y,
		)
		m.Rotate(float64(self.rotation))
		// Then move to the absolute position.
		m.Translate(self.position.X, self.position.Y)

		// And finally move to the origin offset
		if !self.origin.IsZero() {
			m.Translate(-self.origin.X, -self.origin.Y)
		}

		self.matrix = m

		self.dirty = false
	} else {
		m = self.matrix
	}

	if self.parent != nil {
		pm, _ = self.parent.MatrixForParenting()
		m.Concat(pm)
	}

	return m
}
