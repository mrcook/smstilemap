# ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
# ::             WLA DX compiling shell script v1             ::
# ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
# ::  Based on BAT file from smspower.org/maxim/HowToProgram  ::
# ::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::

FILE="splash"

# Cleanup to avoid confusion
rm -f $FILE.o

# Compile
wla-z80 -o $FILE.o $FILE.asm

# Make a linkfile
echo [objects] > linkfile
echo $FILE.o >> linkfile

# Link
wlalink -d -r -v -s linkfile $FILE.sms

# Fixup for eSMS
rm -f $FILE.sms.sym
mv $FILE.sym $FILE.sms.sym

# Cleanup to avoid mess
rm -f linkfile
rm -f $FILE.o
