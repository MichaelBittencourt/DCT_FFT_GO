CC=go build
SRC=pds.go
TARGET=pds
RM=rm

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC)

.PHONY: test clean all

test: $(TARGET)
	@ echo Running dft
	./$< dft < input.txt | ./$< dft inv
	./$< dft < input2.txt | ./$< dft inv
	./$< dft < input3.txt | ./$< dft inv
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
