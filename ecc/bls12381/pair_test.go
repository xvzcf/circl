package bls12381

import (
	"fmt"
	"testing"

	"github.com/cloudflare/circl/ecc/bls12381/ff"
	"github.com/cloudflare/circl/internal/test"
)

func TestProdPair(t *testing.T) {
	const testTimes = 1 << 5
	const N = 3

	listG1 := [N]*G1{}
	listG2 := [N]*G2{}
	listSc := [N]*Scalar{}
	var ePQn, got Gt

	for i := 0; i < testTimes; i++ {
		got.SetOne()
		for j := 0; j < N; j++ {
			listG1[j] = randomG1(t)
			listG2[j] = randomG2(t)
			listSc[j] = &Scalar{}
			listSc[j].Random()

			ePQ := Pair(listG1[j], listG2[j])
			ePQn.Exp(ePQ, listSc[j])
			got.Mul(&got, &ePQn)
		}

		want := ProdPair(listG1[:], listG2[:], listSc[:])

		if !got.IsEqual(want) {
			test.ReportError(t, got, want)
		}
	}
}

func TestPairBilinear(t *testing.T) {
	testTimes := 1 << 5
	for i := 0; i < testTimes; i++ {
		g1 := G1Generator()
		g2 := G2Generator()
		a := &Scalar{}
		b := &Scalar{}
		a.Random()
		b.Random()
		ab := &Scalar{}
		ab.Mul(a, b)
		p := &G1{}
		q := &G2{}
		p.ScalarMult(a, g1)
		q.ScalarMult(b, g2)
		lhs := Pair(p, q)
		tmp := Pair(g1, g2)
		rhs := &Gt{}
		rhs.Exp(tmp, ab)
		if !lhs.IsEqual(rhs) {
			test.ReportError(t, lhs, rhs)
		}
	}
}

func TestPairIdentity(t *testing.T) {
	g1id := &G1{}
	g2id := &G2{}
	g1 := G1Generator()
	g2 := G2Generator()
	g1id.SetIdentity()
	g2id.SetIdentity()
	one := &Gt{}
	one.SetOne()
	ans := Pair(g1id, g2)
	if !ans.IsEqual(one) {
		test.ReportError(t, ans, one)
	}
	ans = Pair(g1, g2id)
	if !ans.IsEqual(one) {
		test.ReportError(t, ans, one)
	}
}

func BenchmarkMiller(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = miller(g1, g2)
	}
}

func BenchmarkFinalExpo(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()
	f := miller(g1, g2)
	g := ff.EasyExponentiation(f)

	b.Run("EasyExp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ff.EasyExponentiation(f)
		}
	})
	b.Run("HardExp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hardExponentiation(g)
		}
	})
	b.Run("FinalExp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = finalExp(f)
		}
	})
}

func BenchmarkPair(b *testing.B) {
	g1 := G1Generator()
	g2 := G2Generator()

	const N = 3
	listG1 := [N]*G1{}
	listG2 := [N]*G2{}
	listExp := [N]*Scalar{}
	for i := 0; i < N; i++ {
		listG1[i] = new(G1)
		listG1[i].Set(g1)
		listG2[i] = new(G2)
		listG2[i].Set(g2)
		listExp[i] = &Scalar{}
		listExp[i].Random()
	}

	b.Run("Pair", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Pair(g1, g2)
		}
	})
	b.Run(fmt.Sprintf("ProdPair%v", N), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ProdPair(listG1[:], listG2[:], listExp[:])
		}
	})
}
