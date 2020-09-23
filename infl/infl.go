package infl

import (
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
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
	ID     string
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
	const t0 = 100
	unit := big.NewInt(100000000)
	f0 := big.NewRat(t0, t0+1)

	var height int

	balances := make(map[string]Balance)

	newBalance := func(id string, i *big.Int) Balance {
		return Balance{
			ID:     id,
			Height: height,
			Amount: cp(i),
		}
	}

	add := func(id string, i *big.Int) {
		balances[id] = newBalance(id, i)
	}

	value := func(b Balance) *big.Int {
		fnum := f0.Num()
		fdenom := f0.Denom()
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

	totalValue := func(m map[string]Balance) *big.Int {
		var total big.Int
		for _, b := range m {
			total.Add(&total, value(b))
		}
		return &total
	}

	xfer := func(to, from string) error {
		fromBalance, ok := balances[from]
		if !ok {
			return fmt.Errorf("no such source balance: %s", from)
		}
		toBalance, ok := balances[to]
		if !ok {
			toBalance = newBalance(to, big.NewInt(0))
		}
		toValue := value(toBalance)
		fromValue := value(fromBalance)

		delete(balances, from)
		delete(balances, to)

		var total big.Int
		total.Add(toValue, fromValue)

		add(to, &total)
		return nil
	}

	for {
		add(uuid.New().String(), unit)
		height++
		var units big.Rat
		units.SetFrac(totalValue(balances), unit)
		v, _ := units.Float64()
		fmt.Printf("height = %d, total value = %.20f; len=%d\n", height, v, len(balances))
		if len(balances) > 100 {
			for i := 0; i < 5; i++ {
				var from, to string
			LOOP:
				for k := range balances {
					switch {
					case from == "":
						from = k
					case to == "":
						to = k
					default:
						break LOOP
					}
				}
				if err := xfer(to, from); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
