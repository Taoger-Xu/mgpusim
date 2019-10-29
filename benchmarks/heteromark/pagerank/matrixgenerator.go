package pagerank

import "math/rand"

type matrixGenerator struct {
	numNode, numConnection   uint32
	xCoords, yCoords         []uint32
	values                   []float32
	positionOccupied         map[uint32]bool
	xCoordIndex, yCoordIndex map[uint32][]uint32
}

func makeMatrixGenerator(numNode, numConnection uint32) matrixGenerator {
	return matrixGenerator{
		numNode:       numNode,
		numConnection: numConnection,
	}
}

func (g matrixGenerator) generateMatrix() csrMatrix {
	g.init()
	g.generateConnections()
	g.normalize()
	m := g.outputCSRFormat()
	return m
}

func (g *matrixGenerator) init() {
	g.xCoords = make([]uint32, 0, g.numConnection)
	g.yCoords = make([]uint32, 0, g.numConnection)
	g.values = make([]float32, 0, g.numConnection)
	g.positionOccupied = make(map[uint32]bool)
	g.xCoordIndex = make(map[uint32][]uint32)
	g.yCoordIndex = make(map[uint32][]uint32)
}

func (g *matrixGenerator) generateConnections() {
	for i := uint32(0); i < g.numConnection; i++ {
		g.generateOneConnection()
	}
}

func (g *matrixGenerator) normalize() {
	for i := uint32(0); i < g.numNode; i++ {
		sum := g.sumColumn(i)
		if sum == 0 {
			continue
		}

		indexes := g.xCoordIndex[i]
		for _, index := range indexes {
			g.values[index] /= sum
		}
	}
}

func (g matrixGenerator) outputCSRFormat() csrMatrix {
	m := csrMatrix{}
	rowOffset := uint32(0)

	for i := uint32(0); i < g.numNode; i++ {
		cols, values := g.selectRowData(i)
		g.sortRowData(cols, values)

		m.rowOffsets = append(m.rowOffsets, rowOffset)
		m.columnNumbers = append(m.columnNumbers, cols...)
		m.values = append(m.values, values...)
		rowOffset += uint32(len(cols))
	}
	m.rowOffsets = append(m.rowOffsets, rowOffset)

	return m
}

func (g matrixGenerator) selectRowData(
	row uint32,
) (
	cols []uint32,
	values []float32,
) {
	indexes := g.yCoordIndex[row]
	for _, index := range indexes {
		cols = append(cols, g.xCoords[index])
		values = append(values, g.values[index])
	}
	return
}

func (g matrixGenerator) sortRowData(cols []uint32, values []float32) {
	for i := 0; i < len(cols); i++ {
		for j := i; j < len(cols); j++ {
			if cols[i] >= cols[j] {
				cols[i], cols[j] = cols[j], cols[i]
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}

func (g matrixGenerator) sumColumn(i uint32) float32 {
	sum := float32(0)
	indexes := g.xCoordIndex[i]
	for _, index := range indexes {
		sum += g.values[index]
	}
	return sum
}

func (g *matrixGenerator) generateOneConnection() {
	x, y := g.generateUnoccupiedPosition()
	v := rand.Float32()
	g.xCoords = append(g.xCoords, x)
	g.yCoords = append(g.yCoords, y)
	g.values = append(g.values, v)

	if _, ok := g.xCoordIndex[x]; !ok {
		g.xCoordIndex[x] = make([]uint32, 0)
	}
	if _, ok := g.yCoordIndex[y]; !ok {
		g.yCoordIndex[y] = make([]uint32, 0)
	}
	g.xCoordIndex[x] = append(g.xCoordIndex[x], uint32(len(g.values)-1))
	g.yCoordIndex[y] = append(g.yCoordIndex[y], uint32(len(g.values)-1))
}

func (g matrixGenerator) generateUnoccupiedPosition() (x, y uint32) {
	for {
		x = uint32(rand.Int()) % g.numNode
		y = uint32(rand.Int()) % g.numNode
		if !g.isPositionOccupied(x, y) {
			g.markPositionOccupied(x, y)
			return
		}
	}
}

func (g matrixGenerator) isPositionOccupied(x, y uint32) bool {
	_, ok := g.positionOccupied[y*g.numNode+x]
	return ok
}

func (g matrixGenerator) markPositionOccupied(x, y uint32) {
	g.positionOccupied[y*g.numNode+x] = true
}
