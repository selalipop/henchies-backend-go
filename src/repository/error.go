package repository

const InvalidPlayerKeyErr InvalidPlayerKeyError = "invalid player key"

type InvalidPlayerKeyError string

func (e InvalidPlayerKeyError) Error() string {
	return string(e)
}

