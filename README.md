# ebi-math

**ebi-math** is a versatile Go package designed for Ebitengine to assist with various mathematical and geometric computations often needed in game development, graphics programming, or any 2D spatial calculations. This library provides a rich set of tools, including vector operations, transformation matrices, random number generation, and more.

## Overview

This library consolidates and enhances functionalities from various sources into a unified, easy-to-use math toolkit. Here's what you can expect:

### Key Features

- **Vector Operations:**
  - 2D Vector (`Vector`) with operations like addition, subtraction, scaling, rotation, normalization, and more.
  - Functions for vector creation, manipulation, and component-wise math operations.

- **Rectangle Manipulations:**
  - `Rectangle` type for handling 2D rectangular areas, including intersection checks, containment, and transformations considering rotation.

- **Transformations:**
  - `Transform` structure for complex transformations, including positions, scales, rotations, and matrix operations for hierarchical transformations.

- **Point Handling:**
  - `Point` type for integer-based coordinate systems, with methods for conversion between float and integer coordinates.

- **Random Number Generation:**
  - `Rand` type for pseudo-random number generation, offering methods for different data types and distributions, including weighted random selection with `RandPicker`.

- **Utility Functions:**
  - Mathematical utilities like linear interpolation (`Lerp`), cubic interpolation, clamping, angle conversions, and more.

### Core Components

- **Vector:** A 2D vector type with methods for geometric calculations.
- **Rectangle:** Supports operations on 2D rectangles, including rotated rectangles.
- **Transform:** A system for managing transformations in a scene graph-like structure.
- **Point:** For handling discrete 2D points.
- **Rand:** Customized random number generation with methods tailored for game logic or simulations.
- **RandPicker:** Allows for weighted random selection, useful in scenarios where outcomes should have varying probabilities.

### Usage Examples

```go
import ebimath "github.com/edwinsyarief/ebi-math"

// Creating a vector
v := ebimath.V(3.0, 4.0)
fmt.Println(v.Length()) // Output: 5.0

// Handling transformations
t := ebimath.T()
t.SetPosition(ebimath.V(10, 10))
t.SetRotation(math.Pi / 4)
fmt.Println(t.Matrix()) // Outputs transformation matrix

// Using random number generation
r := ebimath.Random()
fmt.Println(r.FloatRange(0, 10)) // Random float between 0 and 10

// Weighted random selection
picker := ebimath.RandomPicker[int](r)
picker.AddOption(1, 0.3)
picker.AddOption(2, 0.7)
fmt.Println(picker.Pick()) // Either 1 or 2, with probabilities 30% and 70%
```

### Installation

To use **ebi-math** in your Go project:

```sh
go get github.com/edwinsyarief/ebi-math
```
