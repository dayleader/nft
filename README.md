## DEPRECATED
Please use https://github.com/golang-enthusiast/candy-machine/tree/main

## About 
This repository provides various utilities to help you build your NFT collection:
- Generate images from source layers / Merge layers in the specified order
- Generate ERC-721 traits
- Upload images to IPFS
- Upload metadata to IPFS

## How to generate images?

Run the following command:
`./nft --generated-image-input=${IMG_DIR} --generated-image-output=${OUTPUT_DIR} --generated-image-length=2048 --generated-image-width=2048 --generated-image-number=100 --generate`

### Flags
- generated-image-input - image input directory where the source layers are located (default is input-dir)
- generated-image-output - image output directory  where generated images will be saved here (default is output-dir)
- generated-image-length - canvas length (default is 2048 px)
- generated-image-width - canvas width (default is 2048 px)
- generated-image-number - the number of images to be generated (default is 100 images)

### What format should my layers be in?
To merge layers in the correct order, the input layers directory must be in the following format:
- input-dir:
    - 01_trait_group
        - trait_name.png
        - trait_name.png
        - trait_name.png
    - 02_trait_group
        - trait_name.png
    - 03_trait_group
        - trait_name.png
        - ...

For example:
- input-dir:
    - 01_background
        - ice_cave.png
        - pink_city.png
        - heaven.png
    - 02_body
        - aliens.png
        - crab.png
    - 03_head
        - einstein.png
        - aviator.png
    - 04_pants
        - green_pants.png
        - safari_pants.png
    - 05_shoes
        - black_sneakers.png
        - yellow_sneakers.png
    - 06_accessories
        - red headphones.png
        - orange headphones.png
    - 07_etc ...

### How to make some traits more rare than others?  
Basically all traits are equal, but you can make one trait more rare than others by adding .silver or .gold postfix to your layer.
- .silver - makes a trait 2 times less frequent than others
- .gold - makes a trait 4 times less frequent than others

Example:
Let's make some traits more rare than others. Then the directory of input images should have the following structure:
- input-dir:
    - 01_background
        - ice_cave.gold.png <-- this background will be 4 times less frequent than others
        - pink_city.png
        - heaven.png
    - 02_body
        - aliens.png
        - crab.png
    - 03_head
        - einstein.png
        - aviator.silver.png <-- this background will be 2 times less frequent than others
    - 04_pants
        - green_pants.png
        - safari_pants.png
    - 05_shoes
        - black_sneakers.png
        - yellow_sneakers.png
    - 06_accessories
        - super rare headphones.gold.png <-- this 06_accessorie will be 4 times less frequent than others
        - red headphones.png
        - orange headphones.png
    - 07_etc ...

## How to upload images to IPFS/Pinata?

Run the following command:
`./nft --ipfs-input=${INPUT_DIRECTORY} --ipfs-output=${IPFS_OUTPUT_DIRECTORY} --ipfs-api-key=${PINATA_KEY} --ipfs-secret-key=${PINATA_SECRET} --ipfs-upload`

### Flags
- ipfs-input - input directory / images to upload
- ipfs-output - output directory / ipfs hash & meatadata will be saved here
- ipfs-api-key - Pinata API key
- ipfs-secret-key - Pinata API secret key

## Print statistics
Run the following command to print statistics:
`./nft --generated-image-output=${IMG_DIR} --ipfs-output=${IPFS_OUTPUT_DIRECTORY} --info`
