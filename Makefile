preocv:
	(brew install opencv && brew install pkgconfig) || (cd $${GOPATH}/src/gocv.io/x/gocv && make install)

decode:
	invisible-watermark -v -a decode -t bytes -m rivaGan -w 'REWD' -l 32 ./assets/monet_encoded_gan.png
	invisible-watermark -v -a decode -t bytes -m dwtDctSvd -w 'REWD' -l 32 ./assets/monet_encoded_freq.png

test_burn:
	go test ./test -v

decode_burnt:
	invisible-watermark -v -a decode -t bytes -m dwtDctSvd -w 'REWD' -l 32 ./assets/monet_encoded_freq_burnt.png || true
	invisible-watermark -v -a decode -t bytes -m rivaGan -w 'REWD' -l 32 ./assets/monet_encoded_gan_burnt.png || true

.PHONY: preocv decode test_burn decode_burnt