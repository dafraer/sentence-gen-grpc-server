package currency

type (
	MicroUSD int64
	USD      int64
)

const conversionRate = 1_000_000

func (m MicroUSD) USD() USD {
	return USD(int(m) / conversionRate)
}

func (u USD) MicroUSD() MicroUSD {
	return MicroUSD(int(u) * conversionRate)
}
