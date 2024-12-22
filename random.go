package ebimath

import (
	"math"
	"sort"
	"time"

	"golang.org/x/exp/rand"
)

type Rand struct {
	rng *rand.Rand
}

func Random() *Rand {
	return RandomWidthSeed(time.Now().UnixNano())
}

func RandomWidthSeed(seed int64) *Rand {
	result := &Rand{}
	result.SetSeed(seed)
	return result
}

func (self *Rand) SetSeed(seed int64) {
	self.rng = rand.New(rand.NewSource(uint64(seed)))
}

func (self *Rand) Offset(min, max float64) Vector {
	return Vector{X: self.FloatRange(min, max), Y: self.FloatRange(min, max)}
}

func (self *Rand) Chance(probability float64) bool {
	return self.rng.Float64() <= probability
}

func (self *Rand) Bool() bool {
	return self.rng.Float64() < 0.5
}

func (self *Rand) IntRange(min, max int) int {
	return min + self.rng.Intn(max-min+1)
}

func (self *Rand) PositiveInt64() int64 {
	return self.rng.Int63()
}

func (self *Rand) PositiveInt() int {
	return self.rng.Int()
}

func (self *Rand) Uint64() uint64 {
	return self.rng.Uint64()
}

func (self *Rand) Float64() float64 {
	return self.rng.Float64()
}

func (self *Rand) NextFloat64(max float64) float64 {
	return self.rng.Float64() * max
}

func (self *Rand) FloatRange(min, max float64) float64 {
	return min + self.rng.Float64()*(max-min)
}

func (self *Rand) Rad() float64 {
	return self.FloatRange(0, 2*math.Pi)
}

func (self *Rand) VectorRange(min, max Vector) Vector {
	return min.Add(V(self.NextFloat64(max.X-min.X), self.NextFloat64(max.Y-min.Y)))
}

func RandomIndex[T any](r *Rand, slice []T) int {
	if len(slice) == 0 {
		return -1
	}
	return r.IntRange(0, len(slice)-1)
}

func RandomElement[T any](r *Rand, slice []T) (element T) {
	if len(slice) == 0 {
		return element // Zero value
	}
	if len(slice) == 1 {
		return slice[0]
	}
	return slice[RandomIndex(r, slice)]
}

func RandomChoose[T any](r *Rand, elements ...T) (element T) {
	return RandomElement(r, elements)
}

func RandomShuffle[T any](r *Rand, slice []T) {
	r.rng.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// RandPicker performs a uniformly distributed random probing among the given objects with weights.
// Higher the weight, higher the chance of that object of being picked.
type RandPicker[T any] struct {
	r *Rand

	keys   randPickerKeySlice
	values []T

	threshold float64
	sorted    bool
}

type randPickerKey struct {
	index     int
	threshold float64
}

type randPickerKeySlice []randPickerKey

func (self *randPickerKeySlice) Len() int { return len(*self) }
func (self *randPickerKeySlice) Less(i, j int) bool {
	return (*self)[i].threshold < (*self)[j].threshold
}
func (self *randPickerKeySlice) Swap(i, j int) { (*self)[i], (*self)[j] = (*self)[j], (*self)[i] }

func RandomPicker[T any](r *Rand) *RandPicker[T] {
	return &RandPicker[T]{r: r}
}

func (self *RandPicker[T]) Reset() {
	self.keys = self.keys[:0]
	self.values = self.values[:0]
	self.threshold = 0
	self.sorted = false
}

func (self *RandPicker[T]) AddOption(value T, weight float64) {
	if weight == 0 {
		return // Zero probability in any case
	}
	self.threshold += weight
	self.keys = append(self.keys, randPickerKey{
		threshold: self.threshold,
		index:     len(self.values),
	})
	self.values = append(self.values, value)
	self.sorted = false
}

func (self *RandPicker[T]) AddOptions(values ...T) {
	for _, val := range values {
		self.AddOption(val, 1)
	}
}

func (self *RandPicker[T]) IsEmpty() bool {
	return len(self.values) != 0
}

func (self *RandPicker[T]) Pick() T {
	var result T
	if len(self.values) == 0 {
		return result // Zero value
	}
	if len(self.values) == 1 {
		return self.values[0]
	}

	// In a normal use case the random picker is initialized and then used
	// without adding extra options, so this sorting will happen only once in that case.
	if !self.sorted {
		sort.Sort(&self.keys)
		self.sorted = true
	}

	roll := self.r.FloatRange(0, self.threshold)
	i := sort.Search(len(self.keys), func(i int) bool {
		return roll <= self.keys[i].threshold
	})
	if i < len(self.keys) && roll <= self.keys[i].threshold {
		result = self.values[self.keys[i].index]
	} else {
		result = self.values[len(self.values)-1]
	}
	return result
}
