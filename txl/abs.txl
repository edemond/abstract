comments
		//
end comments

define program
		[repeat statement]
end define

define statement
		[let]
	|	[default]
	| 	[expr]
end define

define let
		'let [id] '= [line]
	| 	'let [id] '= [block]
end define

define default
		'default [expr]
end define

define block
		'{ [repeat line] '}
end define

define line
		[expr]
	|	[expr] [repeat barExpr]
end define

% An expression is a sequence of values that compile to a playable output.
% It's the unit of playability. 
define expr
		[repeat value]
end define

define barExpr
		'| [expr]
end define

% A value is a component of an expression.
define value
		[scale] 
	|	[pitch]
	|	[id] % Variable reference to an expr.
	|   [func]
end define

define func
	[id] '( [list integernumber+] ')
end define

compounds
	'C# 'D# 'E# 'F# 'G# 'A# 'B#
end compounds

define pitch
		'C | 'C# | 'Cb
	|	'D | 'D# | 'Db
	|	'E | 'E# | 'Eb
	|	'F | 'F# | 'Fb
	|	'G | 'G# | 'Gb
	|	'A | 'A# | 'Ab
	|	'B | 'B# | 'Bb
end define

define scale
		'major
	|	'minor
	|	'diminished
	|	'wholetone
	|	'chromatic
	|	'ionian
	|	'dorian
	|	'phrygian
	|	'lydian
	|	'mixolydian
	|	'aeolian
	|	'locrian
end define

% Just copy input to output; we're only testing parsing.
function main
	replace [program]
		P [program]
	by
		P
end function

