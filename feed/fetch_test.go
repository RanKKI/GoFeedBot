package feed

import (
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestLatestTime(t *testing.T) {
    ass := assert.New(t)

    f := Fetcher{}
    t1 := time.Now()
    t2 := t1.Add(10 * time.Second)
    t3, err := f.getLatestTime(&t1, &t2)
    ass.Nil(err)
    ass.Equal(*t3, t2)

    t3, err = f.getLatestTime(&t1, nil)
    ass.Nil(err)
    ass.Equal(*t3, t1)

    t3, err = f.getLatestTime(nil, nil)

    ass.Nil(t3)
    ass.NotNil(err)

}
