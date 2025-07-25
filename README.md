# image_fingerprinting
Jpeg image decoding, fingerprinting & perceptual hashing

## Usage examples:
`go run main.go types.go huffman.go q_tables.go idct.go black.jpeg black_2.jpeg` 

<img width="1366" height="181" alt="Screenshot (83)" src="https://github.com/user-attachments/assets/f5a01fc2-01e1-4f22-88ff-a6bc0d6c1bab" />

## Resources:
https://shorturl.at/onWSx <br>
https://stackoverflow.com/questions/37217640/what-are-the-last-2-bytes-in-the-start-of-scan-of-a-jpeg-jfif-image <br>
https://www.ece.ucf.edu/~mingjie/EEL4783_2012/PA2.pdf <br>
and many more...

## Assumptions:
`4:2:0` chroma subsampling ...<br>
`8 bit` quantization ... <br>
Baseline encoding as opposed to Progressive, Extended, Custom etc ... <br>
Standard table definations ... <br>
That the image is a jpeg ... <br>

## Processing: (view src)
Decode image (huffman encoding through to idct), keep low frequency components of the image (`y` component: `luminance`) effectively turning the image to a gray scale, generate a hash from this array and generate a similarity score based on hamming distance as documented in `main.go`

## Next steps:
Reducing assumptions, unlocking for a wider array of encoding and quantization techniques: `16 bit` (although rare), and better code documenting

## Important:
Similarity is based on features like texture as opposed to colour schemes (considering we are discarding the chrominance components: `cb` `cr`), and the code is format locked for the `jpeg` image format (view jpeg specs)

