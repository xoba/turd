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

type Config struct {
	Height          int
	InflationFactor *big.Rat
	Balances        map[string]Balance
}

func NewConfig(f *big.Rat) Config {
	return Config{
		InflationFactor: f,
		Balances:        make(map[string]Balance),
	}
}

func (c Config) value(b Balance) *big.Int {
	fnum := c.InflationFactor.Num()
	fdenom := c.InflationFactor.Denom()
	dt := big.NewInt(int64(c.Height - b.Height))
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
func (c Config) totalValue() *big.Int {
	var total big.Int
	for _, b := range c.Balances {
		total.Add(&total, c.value(b))
	}
	return &total
}

func (c Config) add(id string, i *big.Int) {
	c.Balances[id] = Balance{
		ID:     id,
		Amount: i,
		Height: c.Height,
	}
}

func (c Config) xfer(to, from string) error {
	fromBalance, ok := c.Balances[from]
	if !ok {
		return fmt.Errorf("no such source balance: %s", from)
	}
	toBalance, ok := c.Balances[to]
	if !ok {
		toBalance = Balance{
			ID:     to,
			Height: c.Height,
			Amount: big.NewInt(0),
		}
	}
	toValue := c.value(toBalance)
	fromValue := c.value(fromBalance)
	delete(c.Balances, from)
	delete(c.Balances, to)
	var total big.Int
	total.Add(toValue, fromValue)
	c.add(to, &total)
	return nil
}

// inflation with big numbers
func RunBig(cnfg.Config) error {
	const (
		t0   = 100
		unit = 100000000
	)
	c := NewConfig(big.NewRat(t0, t0+1))
	for {
		c.add(uuid.New().String(), big.NewInt(unit))
		c.Height++
		var units big.Rat
		units.SetFrac(c.totalValue(), big.NewInt(unit))
		v, _ := units.Float64()
		fmt.Printf("height = %5d, total value = %.20f; len=%d\n", c.Height, v, len(c.Balances))
		if len(c.Balances) > 100 {
			const n = 5
			// choose n balances to transfer from:
			var from []string
		LOOP:
			for k := range c.Balances {
				switch {
				case len(from) < n:
					from = append(from, k)
				default:
					break LOOP
				}
			}
			id := uuid.New().String()
			for _, f := range from {
				if err := c.xfer(id, f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
