package pcm2wav

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPcmBytes2WavBytes(t *testing.T) {

	var err error

	testPcm, err := ioutil.ReadFile("./test.pcm")
	if err != nil {
		assert.Error(t, err)
	}

	expectedWav, err := ioutil.ReadFile("./expected.wav")
	if err != nil {
		assert.Error(t, err)
	}

	resultWav, err := ConvertBytes(testPcm, 1, 16000, 16)
	if err != nil {
		assert.Error(t, err)
	}

	err = ioutil.WriteFile("./result.wav", resultWav, 0666)

	assert.Equal(t, expectedWav, resultWav, "Unexpected wav data")

}
