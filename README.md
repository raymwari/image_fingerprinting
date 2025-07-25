# image_fingerprinting
Jpeg image decoding, fingerprinting & perceptual hashing

## Usage examples:
`go run main.go types.go huffman.go q_tables.go idct.go black.jpeg black_2.jpeg` 

# Demo:
<img width="1366" height="181" alt="Screenshot (83)" src="https://github.com/user-attachments/assets/f5a01fc2-01e1-4f22-88ff-a6bc0d6c1bab" />

## Resources:
https://shorturl.at/onWSx <br>
https://stackoverflow.com/questions/37217640/what-are-the-last-2-bytes-in-the-start-of-scan-of-a-jpeg-jfif-image <br>
https://github.com/shomali11/util/blob/master/xhashes/xhashes.go <br>
https://www.ece.ucf.edu/~mingjie/EEL4783_2012/PA2.pdf <br>
and more...

## Assumptions:
`4:2:0` chroma subsampling. <br>
8 bit quantization. <br>
Baseline encoding. <br>
Standard table definations. <br>

## Processing: (view src)
Decode image (huffman encoding through to idct), keep low frequency components of the image (`y` component: `luminance`) effectively turning the image to a gray scale, generate a hash from this array and compare the hashes

## Next steps:
Reducing assumptions, unlocking for a wider array of encoding and quantization techniques 16 bit (although rare), implementing a similarity score, better commenting, and fixing the code alignment issues

## Important:
Code locked for the `jpeg` image format (and no, changing the file extension doesn't make it a jpeg) <br>

