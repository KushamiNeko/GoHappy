package handler

//type cache struct {
//symbol    string
//frequency data.Frequency

//exstart time.Time
//exend   time.Time

//series *data.TimeSeries
//}

//type cacheStore []*cache

//func (c *cacheStore) Lookup(src data.DataSource, dt time.Time, symbol string, freq data.Frequency) error {

//start, end := p.chartPeriod(dt, freq)

//for _, q := range p.store {
//if q.symbol == symbol && q.frequency == freq {

//if (start.After(q.start) || start.Equal(q.start)) && (end.Before(q.end) || end.Equal(q.end)) {
//p.series = q.quotes
//return nil
//}

//}
//}

//exstart := start.Add(-500 * 24 * time.Hour)
//exend := end.Add(500 * 24 * time.Hour)

//qs, err := src.Read(exstart, exend, symbol, freq)
//if err != nil {
//return err
//}

//qs.TimeSlice(start, end)

//p.store = append(p.store, &quotesCache{symbol: symbol, frequency: freq, quotes: qs, start: exstart, end: exend})
//p.series = qs

//return nil
//}
