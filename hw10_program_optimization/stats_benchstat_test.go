package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

// go test -benchmem -count=10 -bench=BenchmarkForStat //.
func BenchmarkForStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Helper()
		b.StopTimer()

		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(b, err)
		defer r.Close()

		require.Equal(b, 1, len(r.File))

		data, err := r.File[0].Open()
		require.NoError(b, err)

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(b, err)

		// require.Equal(b, expectedBizStat, stat)
	}
}
