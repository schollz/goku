#directories
resources = ./resources

#files
SYLLABLES := $(resources)/cmudict.0.7a

$(SYLLABLES):
	wget -O $(SYLLABLES) http://sourceforge.net/p/cmusphinx/code/11879/tree/trunk/cmudict/cmudict.0.7a?format=raw

download: $(SYLLABLES)

.DEFAULT: download
