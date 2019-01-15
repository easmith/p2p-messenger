package proto

import (
	"os"
	"reflect"
	"testing"
)

func TestCalcSharedSecret(t *testing.T) {
	publicKey1, privateKey1 := CreateKeyExchangePair()

	publicKey2, privateKey2 := CreateKeyExchangePair()

	secret1 := CalcSharedSecret(publicKey1[:], privateKey2[:])
	secret2 := CalcSharedSecret(publicKey2[:], privateKey1[:])

	equal := reflect.DeepEqual(secret1, secret2)
	t.Logf("Secrets are equals? %v", equal)
}

func TestSaveKey(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Can I create new file?",
			args: args{fileName: "auto-test.key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SaveKey(tt.args.fileName)

			if got == nil {
				t.Fatalf("SaveKey() return nil file")
			}

			fileName := got.Name()

			info, e := os.Stat(fileName)
			if e != nil {
				t.Fatalf(e.Error())
			}

			size := info.Size()
			if size != 32 {
				t.Errorf("Incorrect file size: %v", size)
			}
			e = os.Remove(fileName)
			if e != nil {
				t.Fatalf(e.Error())
			}
		})
	}
}
