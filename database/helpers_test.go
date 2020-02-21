package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_ParseCustomType(t *testing.T) {
	tests := map[string]struct {
		testBytes []byte
		validate  func(t *testing.T, actualErr error, actualElems [][]byte)
	}{
		"base path": {
			testBytes: []byte(`()`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 0)
			},
		},
		"alternate path- no strings": {
			testBytes: []byte(`(1,2)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 2)
				expectedElemOne := []byte(`1`)
				expectedElemTwo := []byte(`2`)

				require.Equal(t, expectedElemOne, actualElems[0])
				require.Equal(t, expectedElemTwo, actualElems[1])
			},
		},
		"alternate path- only strings": {
			testBytes: []byte(`("1","2")`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 2)
				expectedElemOne := []byte(`1`)
				expectedElemTwo := []byte(`2`)

				require.Equal(t, expectedElemOne, actualElems[0])
				require.Equal(t, expectedElemTwo, actualElems[1])
			},
		},
		"alternate path- mix with NULL vals": {
			testBytes: []byte(`("1",2,1.1,true,NULL)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 5)
				expectedElemOne := []byte(`1`)
				expectedElemTwo := []byte(`2`)
				expectedElemThree := []byte(`1.1`)
				expectedElemFour := []byte(`true`)

				require.Equal(t, expectedElemOne, actualElems[0])
				require.Equal(t, expectedElemTwo, actualElems[1])
				require.Equal(t, expectedElemThree, actualElems[2])
				require.Equal(t, expectedElemFour, actualElems[3])
				require.Nil(t, actualElems[4])
			},
		},
		"alternate path- with escaped quotes and double double quotes": {
			testBytes: []byte(`("1\" "",","2"",",1.1,true)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 4)
				expectedElemOne := []byte(`1" ",`)
				expectedElemTwo := []byte(`2",`)
				expectedElemThree := []byte(`1.1`)
				expectedElemFour := []byte(`true`)

				require.Equal(t, expectedElemOne, actualElems[0])
				require.Equal(t, expectedElemTwo, actualElems[1])
				require.Equal(t, expectedElemThree, actualElems[2])
				require.Equal(t, expectedElemFour, actualElems[3])
			},
		},
		"alternate path- with new line characters": {
			testBytes: []byte(`("1\\n,","2"",",1.1,true)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 4)
				expectedElemOne := []byte(`1\n,`)
				expectedElemTwo := []byte(`2",`)
				expectedElemThree := []byte(`1.1`)
				expectedElemFour := []byte(`true`)

				require.Equal(t, expectedElemOne, actualElems[0])
				require.Equal(t, expectedElemTwo, actualElems[1])
				require.Equal(t, expectedElemThree, actualElems[2])
				require.Equal(t, expectedElemFour, actualElems[3])
			},
		},
		"alternate path- empty element": {
			testBytes: []byte(`(,,)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.NoError(t, actualErr)
				require.Len(t, actualElems, 3)

				require.Nil(t, actualElems[0])
				require.Nil(t, actualElems[1])
				require.Nil(t, actualElems[2])
			},
		},
		"exceptional path- invalid start to custom type": {
			testBytes: []byte(`;;`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.Error(t, actualErr)
				require.Contains(t, actualErr.Error(), "unexpected ';'")
			},
		},
		//"exceptional path- invalid start to custom type": {
		//	testBytes: []byte(`;;`),
		//	validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
		//		require.Error(t, actualErr)
		//		require.Contains(t, actualErr.Error(), "unexpected ';'")
		//	},
		//},
		"exceptional path- invalid delimeter": {
			testBytes: []byte(`("foobar";)`),
			validate: func(t *testing.T, actualErr error, actualElems [][]byte) {
				require.Error(t, actualErr)
				require.Contains(t, actualErr.Error(), "unexpected ';'")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			elems, err := ParseCustomType(tc.testBytes, []byte(`,`))
			tc.validate(t, err, elems)
		})
	}
}
