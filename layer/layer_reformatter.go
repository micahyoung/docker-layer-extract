package layer

import (
	"archive/tar"
	"io"
	"io/ioutil"
)

type LayerReformatter struct {
}

func NewLayerReformatter() *LayerReformatter {
	return &LayerReformatter{}
}
func (a *LayerReformatter) LayerHasPaxHeaders(layerTarfileReader io.Reader) (hasLayer bool, err error) {
	tarReader := tar.NewReader(layerTarfileReader)

	for {
		var header *tar.Header

		header, err = tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}

		if len(header.PAXRecords) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (a *LayerReformatter) CopyLayerWithoutPaxHeaders(writer io.Writer, reader io.Reader) (err error) {
	tarReader := tar.NewReader(reader)
	tarWriter := tar.NewWriter(writer)

	for {
		var header *tar.Header

		header, err = tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		newHeader := header
		newHeader.PAXRecords = nil           //remove all PAX records
		newHeader.Format = tar.FormatUnknown // let format be redetermined automatically (PAX records only added when needed)

		err = tarWriter.WriteHeader(newHeader)
		if err != nil {
			return err
		}

		content, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return err
		}

		_, err = tarWriter.Write(content)
		if err != nil {
			return err
		}

	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
