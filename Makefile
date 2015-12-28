#directories
resources := ./resources

#files
SYLLABLES  := $(resources)/cmudict.0.7a
THESAUROUS := $(resources)/th_en_US_new.dat

default: $(SYLLABLES) $(THESAUROUS)

$(SYLLABLES):
	mkdir -p $(resources)
	wget -O $(SYLLABLES) http://sourceforge.net/p/cmusphinx/code/11879/tree/trunk/cmudict/cmudict.0.7a?format=raw

$(THESAUROUS):
	mkdir -p $(resources)
	wget -O $(resources)/MyThes-1.zip http://lingucomponent.openoffice.org/MyThes-1.zip
	unzip -d $(resources) $(resources)/MyThes-1.zip
	cp $(resources)/MyThes-1.0/th_en_US_new.dat $(THESAUROUS)
	rm -rf $(resources)/MyThes-1.0 $(resources)/MyThes-1.zip
