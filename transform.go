package ebimath

// Transformer defines the interface for objects that have transforms.
type Transformer interface {
	GetParentTransform() *Transform
	GetTransform() *Transform
}

// Transform represents a transformation in 2D space.
type Transform struct {
	position, scale, offset, origin Vector
	rotation                        float64
	parent                          *Transform
	isDirty                         bool // A single dirty flag for performance caching.
	worldMatrix                     Matrix
	parentMatrix                    Matrix
	parentInverted                  Matrix
}

// T creates a new Transform with default values.
func T() *Transform {
	return &Transform{
		scale:   V2(1),
		isDirty: true,
	}
}

// Methods for Parent Hierarchy
// ----------------------------
// GetParentTransform returns the parent Transform or nil if there is no parent.
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

// GetTransform returns this Transform.
func (self *Transform) GetTransform() *Transform {
	return self
}

// Transformation Properties
// -------------------------
func (self *Transform) Origin() Vector {
	return self.origin
}

func (self *Transform) SetOrigin(origin Vector) {
	self.isDirty = true
	self.origin = origin
}

func (self *Transform) IsDirty() bool {
	// A single, recursive check for dirty state.
	if self.isDirty {
		return true
	}
	if self.parent != nil {
		return self.parent.IsDirty()
	}
	return false
}

// Position and Movement
// ---------------------
// SetPosition updates the position, preserving the world-space position
// by adjusting the local position based on the parent's inverse matrix.
func (self *Transform) SetPosition(position Vector) {
	self.isDirty = true
	if self.parent != nil {
		worldToLocalMatrix := self.parent.Matrix()
		worldToLocalMatrix.Invert()
		self.position = position.Apply(worldToLocalMatrix)
	} else {
		self.position = position
	}
}

// Position returns the absolute position in world space.
func (self *Transform) Position() Vector {
	return V2(0).Apply(self.Matrix())
}

func (self *Transform) Move(v ...Vector) {
	self.SetPosition(self.Position().Add(v...))
}

// Rotation
// --------
// SetRotation updates the rotation, preserving the world-space rotation
// by adjusting the local rotation based on the parent's rotation.
func (self *Transform) SetRotation(rotation float64) {
	self.isDirty = true
	if self.parent != nil {
		self.rotation = rotation - self.parent.Rotation()
	} else {
		self.rotation = rotation
	}
}

// Rotation returns the absolute rotation in world space.
func (self *Transform) Rotation() float64 {
	if self.parent == nil {
		return self.rotation
	}
	return self.rotation + self.parent.Rotation()
}

func (self *Transform) Rotate(rotation float64) {
	self.isDirty = true
	self.rotation += rotation
}

// Scale
// -----
func (self *Transform) SetScale(scale Vector) {
	self.isDirty = true
	self.scale = scale
}

func (self *Transform) Scale() Vector {
	if self.parent == nil {
		return self.scale
	}
	return self.scale.Mul(self.parent.Scale())
}

func (self *Transform) AddScale(add ...Vector) {
	self.isDirty = true
	self.scale = self.scale.Add(add...)
}

// Offset
// ------
func (self *Transform) SetOffset(offset Vector) {
	self.isDirty = true
	self.offset = offset
}

func (self *Transform) Offset() Vector {
	return self.offset
}

// Transform Modifiers
// -------------------
// Abs returns an absolute transform without considering parents.
func (self *Transform) Abs() Transform {
	abs := *T()
	abs.SetPosition(self.Position())
	abs.SetRotation(self.Rotation())
	abs.SetScale(self.Scale())
	abs.SetOffset(self.Offset())
	abs.SetOrigin(self.Origin())
	return abs
}

func (self *Transform) Rel() Transform {
	rel := *self
	rel.parent = nil
	return rel
}

// Parent Management
// -----------------
func (self *Transform) Connected() bool {
	return self.parent != nil
}

// Replace updates this transform's local properties to match the world properties
// of another transform.
func (self *Transform) Replace(new Transformer) {
	nt := new.GetTransform()
	self.SetPosition(nt.Position())
	self.SetOffset(nt.Offset())
	self.SetRotation(nt.Rotation())
	self.SetOrigin(nt.Origin())
	self.SetScale(nt.Scale())
}

// Connect establishes a parent-child relationship, preserving the object's
// world space transform.
func (self *Transform) Connect(parent Transformer) {
	if parent == nil {
		return
	}
	worldPos := self.Position()
	worldRot := self.Rotation()
	worldScale := self.Scale()
	worldOffset := self.Offset()
	worldOrigin := self.Origin()

	self.parent = parent.GetTransform()
	self.isDirty = true

	self.SetPosition(worldPos)
	self.SetRotation(worldRot)
	self.SetScale(worldScale)
	self.SetOffset(worldOffset)
	self.SetOrigin(worldOrigin)
}

// Disconnect removes the parent relationship, making the transform absolute.
func (self *Transform) Disconnect() {
	if self.parent == nil {
		return
	}
	*self = self.Abs()
}

// Matrix Operations
// -----------------
// MatrixForParenting returns matrices for child positioning.
// It returns the world matrix without the origin offset and its inverse.
func (self *Transform) MatrixForParenting() (Matrix, Matrix) {
	if self.isDirty {
		self.Matrix() // Ensure the world matrix is computed and cached.
	}
	return self.parentMatrix, self.parentInverted
}

// Matrix computes the full world transformation matrix for this node.
// The result is cached to avoid repeated calculations.
func (self *Transform) Matrix() Matrix {
	if !self.IsDirty() {
		return self.worldMatrix
	}

	// Calculate the local transformation matrix.
	localMatrix := Matrix{}
	localMatrix.Scale(self.scale.X, self.scale.Y)
	localMatrix.Rotate(self.rotation)
	localMatrix.Translate(self.position.X, self.position.Y)

	// Save this local matrix as the parent matrix for children, without the origin offset.
	self.parentMatrix = localMatrix
	self.parentInverted = localMatrix
	self.parentInverted.Invert()

	// Apply origin and offset translations last, as they should be relative to
	// the object's local space.
	if !self.origin.IsZero() {
		localMatrix.Translate(-self.origin.X, -self.origin.Y)
	}
	if !self.offset.IsZero() {
		localMatrix.Translate(-self.offset.X*self.scale.X, -self.offset.Y*self.scale.Y)
	}

	// If there's a parent, combine this local matrix with the parent's world matrix.
	if self.parent != nil {
		parentMatrix := self.parent.Matrix()
		localMatrix.Concat(parentMatrix)
	}

	self.worldMatrix = localMatrix
	self.isDirty = false
	return self.worldMatrix
}
