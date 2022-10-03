CC=go build
SRC=pds.go
TARGET=pds
RM=rm

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC)

.PHONY: test clean all

test: $(TARGET)
	@ echo Running dct
	./$< dct < input.txt | ./$< dct inv
	./$< dct < input2.txt | ./$< dct inv
	./$< dct < input3.txt | ./$< dct inv
	@ echo Running fft with time decimation
	./$< < input.txt | ./$< inv
	./$< < input2.txt | ./$< inv
	./$< < input3.txt | ./$< inv
	@ echo Running fft with frequency decimation
	./$< freq < input.txt | ./$< freq inv
	./$< freq < input2.txt | ./$< freq inv
	./$< freq < input3.txt | ./$< freq inv

clean:
	$(RM) $(TARGET) 
