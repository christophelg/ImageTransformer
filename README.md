# ImageTransformer
Tool to apply an image transformation to a directory of images

## Motivation

# Examples

	imagetransformer --sourceDir /a/b/c --targetDir /x/y/z --transformation resize --param 4

This example will:
* scan the sourceDir /a/b/c for *.jpg pictures
* resize the pictures found by divising their size by 4
* recreate the directory structure under /a/b/c under /x/y/z with the resized images


	imagetransformer --sourceDir /a/b/c --targetDir /x/y/z --transformation crop --param 400

This example will:
* scan the sourceDir /a/b/c for *.jpg pictures
* crop the pictures found in a 400x400 rectangle
* recreate the directory structure under /a/b/c under /x/y/z with the resized images
