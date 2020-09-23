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
	Amount *big.Rat
}

// inflation with big numbers
func RunBig(c cnfg.Config) error {
	var height int
	var list []Balance
	add := func(i *big.Rat) {
		list = append(list, Balance{
			Height: height,
			Amount: i,
		})
	}
	const t0 = 100

	//	unit := big.NewInt(1000000)
	f0 := big.NewRat(t0, t0+1)
	fnum := f0.Num()
	fdenom := f0.Denom()

	value := func(b Balance) *big.Rat {
		dt := big.NewInt(int64(height - b.Height))
		var fnum2, fdenom2 big.Int
		fnum2.Exp(fnum, dt, nil)
		fdenom2.Exp(fdenom, dt, nil)
		var f big.Rat
		f.SetFrac(&fnum2, &fdenom2)
		var b2 big.Rat
		b2.Mul(b.Amount, &f)
		return &b2
	}

	totalValue := func() *big.Rat {
		var total big.Rat
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
		add(big.NewRat(1, 1))
		height++
		v, _ := totalValue().Float64()
		fmt.Printf("height = %d, total value = %.20f\n", height, v)
	}

	return nil
}
