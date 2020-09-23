package infl

import (
	"fmt"
	"math/big"
	"time"

	"github.com/xoba/turd/cnfg"
)

// study inflation methods, like exponential mining or decaying value
func Run(c cnfg.Config) error {
	last := time.Now()
	var circulating, produced float64
	const t0 = 300
	f := float64(t0) / float64(t0+1)
	fmt.Printf("f = %f\n", f)
	for {
		produced++
		circulating += 1
		circulating *= f
		if time.Since(last) > 100*time.Millisecond {
			fmt.Printf("total = %f / %.2f (%f)\n",
				circulating, produced/t0, produced,
			)
			last = time.Now()
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}

type Balance struct {
	Height int
	Amount *big.Int
}

func cp(i *big.Int) *big.Int {
	var x big.Int
	return x.Set(i)
}

func approx(r *big.Rat) *big.Int {
	var z big.Int
	z.Div(r.Num(), r.Denom())
	return &z
}

// inflation with big numbers
func RunBig(c cnfg.Config) error {
	var height int
	var list []Balance
	add := func(i *big.Int) {
		list = append(list, Balance{
			Height: height,
			Amount: cp(i),
		})
	}
	const t0 = 100

	unit := big.NewInt(100000000)
	f0 := big.NewRat(t0, t0+1)
	fnum := f0.Num()
	fdenom := f0.Denom()

	value := func(b Balance) *big.Int {
		dt := big.NewInt(int64(height - b.Height))
		var fnum2, fdenom2 big.Int
		fnum2.Exp(fnum, dt, nil)
		fdenom2.Exp(fdenom, dt, nil)
		var f big.Rat
		f.SetFrac(&fnum2, &fdenom2)
		var b0 big.Rat
		b0.SetFrac(b.Amount, big.NewInt(1))
		var b2 big.Rat
		b2.Mul(&b0, &f)
		return approx(&b2)
	}

	totalValue := func() *big.Int {
		var total big.Int
		for _, b := range list {
			total.Add(&total, value(b))
		}
		if len(list) > 10 {
			// consolidate list:
			list = list[:0]
			add(&total)
		}
		return &total
	}

	for {
		add(unit)
		height++
		var units big.Rat
		units.SetFrac(totalValue(), unit)
		v, _ := units.Float64()
		fmt.Printf("height = %d, total value = %.20f\n", height, v)
		time.Sleep(time.Millisecond)
	}

	return nil
}
