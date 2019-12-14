package elektron

type Pattern struct {
	Bank, Pattern int
}

func (p Pattern) Program() uint8 {
	return uint8(p.Bank*16 + p.Pattern)
}

var (
	A01 = Pattern{0, 0}
	A02 = Pattern{0, 1}
	A03 = Pattern{0, 2}
	A04 = Pattern{0, 3}
	A05 = Pattern{0, 4}
	A06 = Pattern{0, 5}
	A07 = Pattern{0, 6}
	A08 = Pattern{0, 7}
	A09 = Pattern{0, 8}
	A10 = Pattern{0, 9}
	A11 = Pattern{0, 10}
	A12 = Pattern{0, 11}
	A13 = Pattern{0, 12}
	A14 = Pattern{0, 13}
	A15 = Pattern{0, 14}
	A16 = Pattern{0, 15}

	B01 = Pattern{1, 0}
	B02 = Pattern{1, 1}
	B03 = Pattern{1, 2}
	B04 = Pattern{1, 3}
	B05 = Pattern{1, 4}
	B06 = Pattern{1, 5}
	B07 = Pattern{1, 6}
	B08 = Pattern{1, 7}
	B09 = Pattern{1, 8}
	B10 = Pattern{1, 9}
	B11 = Pattern{1, 10}
	B12 = Pattern{1, 11}
	B13 = Pattern{1, 12}
	B14 = Pattern{1, 13}
	B15 = Pattern{1, 14}
	B16 = Pattern{1, 15}

	C01 = Pattern{2, 0}
	C02 = Pattern{2, 1}
	C03 = Pattern{2, 2}
	C04 = Pattern{2, 3}
	C05 = Pattern{2, 4}
	C06 = Pattern{2, 5}
	C07 = Pattern{2, 6}
	C08 = Pattern{2, 7}
	C09 = Pattern{2, 8}
	C10 = Pattern{2, 9}
	C11 = Pattern{2, 10}
	C12 = Pattern{2, 11}
	C13 = Pattern{2, 12}
	C14 = Pattern{2, 13}
	C15 = Pattern{2, 14}
	C16 = Pattern{2, 15}

	D01 = Pattern{3, 0}
	D02 = Pattern{3, 1}
	D03 = Pattern{3, 2}
	D04 = Pattern{3, 3}
	D05 = Pattern{3, 4}
	D06 = Pattern{3, 5}
	D07 = Pattern{3, 6}
	D08 = Pattern{3, 7}
	D09 = Pattern{3, 8}
	D10 = Pattern{3, 9}
	D11 = Pattern{3, 10}
	D12 = Pattern{3, 11}
	D13 = Pattern{3, 12}
	D14 = Pattern{3, 13}
	D15 = Pattern{3, 14}
	D16 = Pattern{3, 15}

	E01 = Pattern{4, 0}
	E02 = Pattern{4, 1}
	E03 = Pattern{4, 2}
	E04 = Pattern{4, 3}
	E05 = Pattern{4, 4}
	E06 = Pattern{4, 5}
	E07 = Pattern{4, 6}
	E08 = Pattern{4, 7}
	E09 = Pattern{4, 8}
	E10 = Pattern{4, 9}
	E11 = Pattern{4, 10}
	E12 = Pattern{4, 11}
	E13 = Pattern{4, 12}
	E14 = Pattern{4, 13}
	E15 = Pattern{4, 14}
	E16 = Pattern{4, 15}

	F01 = Pattern{5, 0}
	F02 = Pattern{5, 1}
	F03 = Pattern{5, 2}
	F04 = Pattern{5, 3}
	F05 = Pattern{5, 4}
	F06 = Pattern{5, 5}
	F07 = Pattern{5, 6}
	F08 = Pattern{5, 7}
	F09 = Pattern{5, 8}
	F10 = Pattern{5, 9}
	F11 = Pattern{5, 10}
	F12 = Pattern{5, 11}
	F13 = Pattern{5, 12}
	F14 = Pattern{5, 13}
	F15 = Pattern{5, 14}
	F16 = Pattern{5, 15}

	G01 = Pattern{6, 0}
	G02 = Pattern{6, 1}
	G03 = Pattern{6, 2}
	G04 = Pattern{6, 3}
	G05 = Pattern{6, 4}
	G06 = Pattern{6, 5}
	G07 = Pattern{6, 6}
	G08 = Pattern{6, 7}
	G09 = Pattern{6, 8}
	G10 = Pattern{6, 9}
	G11 = Pattern{6, 10}
	G12 = Pattern{6, 11}
	G13 = Pattern{6, 12}
	G14 = Pattern{6, 13}
	G15 = Pattern{6, 14}
	G16 = Pattern{6, 15}

	H01 = Pattern{7, 0}
	H02 = Pattern{7, 1}
	H03 = Pattern{7, 2}
	H04 = Pattern{7, 3}
	H05 = Pattern{7, 4}
	H06 = Pattern{7, 5}
	H07 = Pattern{7, 6}
	H08 = Pattern{7, 7}
	H09 = Pattern{7, 8}
	H10 = Pattern{7, 9}
	H11 = Pattern{7, 10}
	H12 = Pattern{7, 11}
	H13 = Pattern{7, 12}
	H14 = Pattern{7, 13}
	H15 = Pattern{7, 14}
	H16 = Pattern{7, 15}
)
