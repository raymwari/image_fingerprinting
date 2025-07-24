# image_fingerprinting
Image (.jpeg) fingerprinting in my fav toy language

## Usage examples:
`go run main.go types.go huffman.go q_tables.go idct.go black.jpeg black_2.jpeg` or <br>
`go build main.go types.go huffman.go q_tables.go idct.go black.jpeg black_2.jpeg` then <br>
`main.exe black.jpeg black_2.jpeg`

## Resources:
https://forensics.map-base.info/report_1/index_en.shtml#:~:text=JPG%20images%20always%20start%20with,using%20byte%20code%20%22FFD9%22. <br>
https://stackoverflow.com/questions/37217640/what-are-the-last-2-bytes-in-the-start-of-scan-of-a-jpeg-jfif-image <br>
https://github.com/shomali11/util/blob/master/xhashes/xhashes.go <br>
https://www.ece.ucf.edu/~mingjie/EEL4783_2012/PA2.pdf <br>
and more...

## Assumptions:
`4:2:0` chroma subsampling <br>
8 bit quantization <br>
Baseline encoding <br>
Standard table definations <br>

## Processing: (view src)
1. Decode jpeg <br>
2. Keep low frequency components (luminance) <br>
3. Generate hash <br>
4. Compare

## Future steps:
Reduce assumptions <br>
Implement a similarity score (hamming distance) <br>
Better comments <br>
Clean code <br>

