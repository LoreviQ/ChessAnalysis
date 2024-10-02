package main

import (
	"github.com/LoreviQ/ChessAnalysis/app_go/internal/database"
)

func main() {
	d := database.NewConnection()
	moves_to_insert := []string{
		"e4", "e5", "Nc3", "Nc6", "f4", "exf4", "Nf3", "Bb4",
		"d4", "Bxc3+", "bxc3", "d5", "e5", "f6", "Bxf4", "fxe5",
		"Bxe5", "Nxe5", "Nxe5", "Qe7", "Bd3", "c5", "O-O", "Nf6",
		"Qf3", "Bg4", "Qf2", "O-O", "Qg3", "Ne4", "Bxe4", "dxe4",
		"Qxg4", "cxd4", "cxd4", "Rad8", "c3", "b5", "Qxe4", "Qa3",
		"Rf3", "a5", "Raf1", "Qxa2", "Rxf8+", "Rxf8", "Rxf8+",
		"Kxf8", "Qa8+", "Ke7", "Nc6+", "Ke6", "Nxa5", "Qe2",
		"Qc6+", "Kf5", "Qf3+", "Qxf3", "gxf3", "Kf4", "Kf2",
		"g5", "d5", "Ke5", "Nc6+", "Kxd5", "Nd4", "b4", "cxb4",
		"Kxd4", "Kg3", "Kc4", "Kg4", "h6", "Kh5", "Kxb4",
		"Kxh6", "Kc5", "Kxg5", "Kd6", "Kg6", "Ke5", "h4",
		"Kf4", "h5", "Kxf3", "h6", "Ke2", "h7", "Kd3", "h8=Q",
		"Kc4", "Qe5", "Kd3", "Kf5", "Kc4", "Kf4", "Kd3",
		"Qe4+", "Kc3", "Kf3", "Kd2", "Qe3+", "Kc2", "Kf2",
		"Kd1", "Qe2+", "Kc1", "Ke3", "Kb1", "Kd3", "Kc1",
		"Qf2", "Kd1", "Qf1#",
	}
	d.InsertMoves(moves_to_insert, "123456")
	d.Close()
}
