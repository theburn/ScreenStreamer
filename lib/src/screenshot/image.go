package screenshot

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

struct YCbCr
{
	unsigned char Y;
	unsigned char Cb;
	unsigned char Cr;
};

struct YCbCr RGBToYCbCr(unsigned char r, unsigned char g, unsigned char b) {
	float fr = (float)r / 255;
	float fg = (float)g / 255;
	float fb = (float)b / 255;

	struct YCbCr ycbcr;
	ycbcr.Y = (unsigned char)((0.2990 * fr + 0.5870 * fg + 0.1140 * fb) * 255);
	ycbcr.Cb = (unsigned char)((-0.1687 * fr - 0.3313 * fg + 0.5000 * fb) * 255 + 128);
	ycbcr.Cr = (unsigned char)((0.5000 * fr - 0.4187 * fg - 0.0813 * fb) * 255 + 128);

	return ycbcr;
}

struct YCbCr RGBToYCbCr4(unsigned char r, unsigned char g, unsigned char b) {
	struct YCbCr ycbcr;
	ycbcr.Y = (unsigned char)((19595 * r + 38470 * g + 7471 * b) >> 16);
	ycbcr.Cb = (unsigned char)(((-11056 * r - 21712 * g + 32768 * b) >> 16) + 128);
	ycbcr.Cr = (unsigned char)(((32768 * r - 27440 * g - 5328 * b) >> 16) + 128);

	return ycbcr;
}

struct YCbCr RGBToYCbCr3(unsigned char r, unsigned char g, unsigned char b) {
	struct YCbCr ycbcr;
	ycbcr.Y = (unsigned char)((2990 * r + 5870 * g + 1140 * b) / 10000);
	ycbcr.Cb = (unsigned char)((-1687 * r - 3313 * g + 5000 * b) / 10000  + 128);
	ycbcr.Cr = (unsigned char)((5000 * r - 4187 * g - 813 * b) / 10000  + 128);

	return ycbcr;
}

static struct YCbCr RGBToYCbCr2(unsigned char r, unsigned char g, unsigned char b) {
	float fr = (float)r / 255;
	float fg = (float)g / 255;
	float fb = (float)b / 255;

	struct YCbCr ycbcr;
	int Y = (int)((0.2990 * fr + 0.5870 * fg + 0.1140 * fb) * 255);
	int Cb = (int)((-0.1687 * fr - 0.3313 * fg + 0.5000 * fb) * 255 + 128);
	int Cr = (int)((0.5000 * fr - 0.4187 * fg - 0.0813 * fb) * 255 + 128);
	if (Y < 0) {
		Y = 0;
	} else if (Y > 0xff) {
		Y = 0xff;
	}
	if (Cb < 0) {
		Cb = 0;
	} else if (Cb > 0xff) {
		Cb = 0xff;
	}
	if (Cr < 0) {
		Cr = 0;
	} else if (Cr > 0xff) {
		Cr = 0xff;
	}
	ycbcr.Y = (unsigned char)(Y);
	ycbcr.Cb = (unsigned char)(Cb);
	ycbcr.Cr = (unsigned char)(Cr);

	return ycbcr;
}

static void ImageRGBToYCbCr444(unsigned char *data, int32_t length, unsigned char *y, unsigned char *cb, unsigned char *cr) {
	int32_t n = 0;
	struct YCbCr ycbcr;
	for (int32_t i = 0; i < length; i += 4) {
		ycbcr = RGBToYCbCr4(data[i+2], data[i+1], data[i]);
		y[n] = ycbcr.Y;
		cb[n] = ycbcr.Cb;
		cr[n] = ycbcr.Cr;
		n += 1;
	}
}

*/
import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

func RGBToYCbCr(r, g, b uint8) (y, cb, cr uint8) {
	ret := C.RGBToYCbCr((C.uchar)(r), (C.uchar)(g), (C.uchar)(b))
	return uint8(ret.Y), uint8(ret.Cb), uint8(ret.Cr)
}

func CRGBToYCbCr444Linux(data, y, cb, cr []byte) {
	C.ImageRGBToYCbCr444((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&y[0])),
		(*C.uchar)(unsafe.Pointer(&cb[0])),
		(*C.uchar)(unsafe.Pointer(&cr[0])))
}

func CRGBToYCbCr444Windows(data, y, cb, cr []byte) {
	C.ImageRGBToYCbCr444((*C.uchar)(unsafe.Pointer(&data[0])),
		C.int32_t(len(data)),
		(*C.uchar)(unsafe.Pointer(&y[0])),
		(*C.uchar)(unsafe.Pointer(&cb[0])),
		(*C.uchar)(unsafe.Pointer(&cr[0])))
}

func RGBAToYCbCr444(img *image.RGBA) *image.YCbCr {
	new_img := image.NewYCbCr(img.Rect, image.YCbCrSubsampleRatio444)
	new_img.Y = make([]uint8, len(img.Pix)/4)
	new_img.Cb = make([]uint8, len(img.Pix)/4)
	new_img.Cr = make([]uint8, len(img.Pix)/4)

	n := 0
	for i := 0; i < len(img.Pix); i += 4 {
		y, cb, cr := color.RGBToYCbCr(img.Pix[i], img.Pix[i+1], img.Pix[i+2])
		new_img.Y[n] = y
		new_img.Cb[n] = cb
		new_img.Cr[n] = cr
		n += 1
	}

	return new_img
}

func RGBAToYCbCr420(img *image.RGBA) *image.YCbCr {
	new_img := image.NewYCbCr(img.Rect, image.YCbCrSubsampleRatio420)
	new_img.Y = make([]uint8, len(img.Pix)/4)
	new_img.Cb = make([]uint8, len(img.Pix)/16)
	new_img.Cr = make([]uint8, len(img.Pix)/16)

	cn := 0
	for y := 0; y < img.Rect.Dy()/2; y += 1 {
		for x := 0; x < img.Rect.Dx()/2; x += 1 {
			x0, y0 := x*2, y*2
			x1, y1 := x*2+1, y*2
			x2, y2 := x*2, y*2+1
			x3, y3 := x*2+1, y*2+1

			co0 := img.RGBAAt(x0, y0)
			cy0, cb0, cr0 := color.RGBToYCbCr(co0.R, co0.G, co0.B)
			co1 := img.RGBAAt(x1, y1)
			cy1, cb1, cr1 := color.RGBToYCbCr(co1.R, co1.G, co1.B)
			co2 := img.RGBAAt(x2, y2)
			cy2, cb2, cr2 := color.RGBToYCbCr(co2.R, co2.G, co2.B)
			co3 := img.RGBAAt(x3, y3)
			cy3, cb3, cr3 := color.RGBToYCbCr(co3.R, co3.G, co3.B)

			new_img.Y[x0+y0*img.Rect.Dx()] = cy0
			new_img.Y[x1+y1*img.Rect.Dx()] = cy1
			new_img.Y[x2+y2*img.Rect.Dx()] = cy2
			new_img.Y[x3+y3*img.Rect.Dx()] = cy3
			new_img.Cb[cn] = cb0/4 + cb1/4 + cb2/4 + cb3/4
			new_img.Cr[cn] = cr0/4 + cr1/4 + cr2/4 + cr3/4
			cn += 1
		}
	}

	return new_img
}