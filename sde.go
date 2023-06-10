package sde

import (
	"bufio"
	"hash"
	"hash/fnv"
	"log"
	"os"
)

type SDE struct {
	hash hash.Hash64

	dictFile *os.File

	curOff int64

	index     map[uint64]int64
	indexFile *os.File
}

func NewSDE(dictFileName, indexFileName string) (*SDE, error) {
	indexFile, err := os.Create(indexFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	dictFile, err := os.Create(dictFileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &SDE{
		hash:      fnv.New64a(),
		dictFile:  dictFile,
		curOff:    0,
		index:     make(map[uint64]int64),
		indexFile: indexFile,
	}, nil
}

func (s *SDE) write(bs []byte, offset int64) error {
	_, err := s.dictFile.WriteAt(bs, int64(offset))
	return err
}

func (s *SDE) Encode(str string) (int64, error) {
	bs := []byte(str)
	bs = append(bs, 0)
	s.hash.Write(bs)
	key := s.hash.Sum64()
	for {
		if off, has := s.index[key]; !has {
			s.index[key] = s.curOff
			if err := s.write(bs, s.curOff); err != nil {
				return -1, err
			}
			preOff := s.curOff
			s.curOff += int64(len(bs))
			return preOff, nil
		} else {
			val, err := s.Decode(off)
			if err != nil {
				return -1, err
			}
			if str == val {
				return off, nil
			}
		}
		key += 1
	}
}

func (s *SDE) Decode(off int64) (string, error) {
	s.dictFile.Seek(off, 0)
	bufReader := bufio.NewReader(s.dictFile)
	bs := make([]byte, 0)
	for {
		b, err := bufReader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		bs = append(bs, b)
	}
	return string(bs), nil
}
