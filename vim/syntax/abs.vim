" Vim syntax file
" Language: Abstract
" Maintainer: Evan DeMond
" Latest Revision: 9 April 2016

if exists("b:current_syntax")
  finish
endif

" Keywords
syn keyword abstractKeyword let default
syn keyword abstractKeyword poly match cutoff
syn keyword abstractKeyword bpm ppq
syn keyword abstractKeyword chord dynamics instrument meter note pitch prob scale voicing 

" Scales
syn keyword abstractBuiltIn major minor 
syn keyword abstractBuiltIn ionian dorian phrygian lydian mixolydian aeolian locrian

" Pitches
syn keyword abstractBuiltIn C C\# C\#\# Cb Cbb 
syn keyword abstractBuiltIn D D\# D\#\# Db Dbb 
syn keyword abstractBuiltIn E E\# E\#\# Eb Ebb 
syn keyword abstractBuiltIn F F\# C\#\# Fb Fbb 
syn keyword abstractBuiltIn G G\# C\#\# Gb Gbb 
syn keyword abstractBuiltIn A A\# C\#\# Ab Abb 
syn keyword abstractBuiltIn B B\# C\#\# Bb Bbb 
syn keyword abstractBuiltIn ppp pp p mp mf f ff fff

" Chords
syn keyword abstractBuiltIn I i II ii III iii IV iv V v VI vi VII vii
syn keyword abstractBuiltIn Ⅰ Ⅱ Ⅲ Ⅳ Ⅴ Ⅵ Ⅶ 
syn keyword abstractBuiltIn ⅰ ⅱ ⅲ ⅳ ⅴ ⅵ ⅶ 
syn match abstractBuiltIn "@[a-zA-Z0-9]*"

" Octaves
syn match abstractBuiltIn "O[a-zA-Z0-9]"

" Comments
syn match abstractComment "//.*$"
syn region abstractComment start="/\*" end="\*/"

" Strings
syn match abstractString "\".*\""
" syn match abstractNumber "\d\+" " durh this doesn't 100% work

" Highlighting rules
hi def link abstractKeyword Statement
hi def link abstractBuiltIn Constant
hi def link abstractString  Constant
hi def link abstractNumber  Constant
hi def link abstractComment Comment
