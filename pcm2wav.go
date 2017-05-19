package pcm2wav

import (
	"encoding/binary"
	"errors"
)

// http://soundfile.sapp.org/doc/WaveFormat/

func littleEndianIntToHex(integer int, numberOfBytes int) (bytes []byte) {
	bytes = make([]byte, numberOfBytes)
	switch numberOfBytes {
	case 2:
		binary.LittleEndian.PutUint16(bytes, uint16(integer))
	case 4:
		binary.LittleEndian.PutUint32(bytes, uint32(integer))
	}
	return
}

func applyString(dst []byte, s string, numberOfBytes int) {
	copy(dst, []byte(s)[:numberOfBytes])
}

func applyLittleEndianInteger(dst []byte, i int, numberOfBytes int) {
	copy(dst, littleEndianIntToHex(i, numberOfBytes)[0:numberOfBytes])
}

type riffChunk struct {
	ChunkId   [4]byte
	ChunkSize [4]byte
	Format    [4]byte
}

func (rc *riffChunk) applyChunkId(chunkId string) {
	applyString(rc.ChunkId[:], chunkId, 4)
}

func (rc *riffChunk) applyChunkSize(chunkSize int) {
	applyLittleEndianInteger(rc.ChunkSize[:], chunkSize, 4)
}

func (rc *riffChunk) applyFormat(format string) {
	applyString(rc.Format[:], format, 4)
}

type fmtSubChunk struct {
	Subchunk1Id   [4]byte
	Subchunk1Size [4]byte
	AudioFormat   [2]byte
	NumChannels   [2]byte
	SampleRate    [4]byte
	ByteRate      [4]byte
	BlockAlign    [2]byte
	BitsPerSample [2]byte
}

func (c *fmtSubChunk) applySubchunk1Id(subchunk1Id string) {
	applyString(c.Subchunk1Id[:], subchunk1Id, 4)
}

func (c *fmtSubChunk) applySubchunk1Size(subchunk1Size int) {
	applyLittleEndianInteger(c.Subchunk1Size[:], subchunk1Size, 4)
}

func (c *fmtSubChunk) applyAudioFormat(audioFormat int) {
	applyLittleEndianInteger(c.AudioFormat[:], audioFormat, 2)
}

func (c *fmtSubChunk) applyNumChannels(numChannels int) {
	applyLittleEndianInteger(c.NumChannels[:], numChannels, 2)
}

func (c *fmtSubChunk) applySampleRate(sampleRate int) {
	applyLittleEndianInteger(c.SampleRate[:], sampleRate, 4)
}

func (c *fmtSubChunk) applyByteRate(byteRate int) {
	applyLittleEndianInteger(c.ByteRate[:], byteRate, 4)
}

func (c *fmtSubChunk) applyBlockAlign(blockAlign int) {
	applyLittleEndianInteger(c.BlockAlign[:], blockAlign, 2)
}

func (c *fmtSubChunk) applyBitsPerSample(bitsPerSample int) {
	applyLittleEndianInteger(c.BitsPerSample[:], bitsPerSample, 2)
}

type dataSubChunk struct {
	Subchunk2Id   [4]byte
	Subchunk2Size [4]byte
}

func (c *dataSubChunk) applySubchunk2Id(subchunk2Id string) {
	applyString(c.Subchunk2Id[:], subchunk2Id, 4)
}

func (c *dataSubChunk) applySubchunk2Size(subchunk2Size int) {
	applyLittleEndianInteger(c.Subchunk2Size[:], subchunk2Size, 4)
}

func PcmBytes2WavBytes(pcm []byte, channels int, sampleRate int, bitsPerSample int) (wav []byte, err error) {
	if channels != 1 && channels != 2 {
		return wav, errors.New("invalid_channels_value")
	}
	if sampleRate != 8000 && sampleRate != 16000 {
		return wav, errors.New("invalid_sample_rate_value")
	}
	if bitsPerSample != 8 && bitsPerSample != 16 {
		return wav, errors.New("invalid_bits_per_sample_value")
	}

	pcmLength := len(pcm)
	subchunk1Size := 16
	subchunk2Size := pcmLength
	chunkSize := 4 + (8 + subchunk1Size) + (8 + subchunk2Size)

	rc := riffChunk{}
	rc.applyChunkId("RIFF")
	rc.applyChunkSize(chunkSize)
	rc.applyFormat("WAVE")

	fsc := fmtSubChunk{}
	fsc.applySubchunk1Id("fmt ")
	fsc.applySubchunk1Size(subchunk1Size)
	fsc.applyAudioFormat(1)
	fsc.applyNumChannels(channels)
	fsc.applySampleRate(sampleRate)
	fsc.applyByteRate(sampleRate * channels * bitsPerSample / 8)
	fsc.applyBlockAlign(channels * bitsPerSample / 8)
	fsc.applyBitsPerSample(bitsPerSample)

	dsc := dataSubChunk{}
	dsc.applySubchunk2Id("data")
	dsc.applySubchunk2Size(subchunk2Size)

	wav = make([]byte, 0, pcmLength+8)

	wav = append(wav, rc.ChunkId[:]...)
	wav = append(wav, rc.ChunkSize[:]...)
	wav = append(wav, rc.Format[:]...)

	wav = append(wav, fsc.Subchunk1Id[:]...)
	wav = append(wav, fsc.Subchunk1Size[:]...)
	wav = append(wav, fsc.AudioFormat[:]...)
	wav = append(wav, fsc.NumChannels[:]...)
	wav = append(wav, fsc.SampleRate[:]...)
	wav = append(wav, fsc.ByteRate[:]...)
	wav = append(wav, fsc.BlockAlign[:]...)
	wav = append(wav, fsc.BitsPerSample[:]...)

	wav = append(wav, dsc.Subchunk2Id[:]...)
	wav = append(wav, dsc.Subchunk2Size[:]...)
	wav = append(wav, pcm...)

	return wav, err
}
