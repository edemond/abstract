# Abstract

Circa 2014. This is an attempt to build a language and interpreter for music composition. The core idea was that the abstractive power of programming languages could be a tool for more quickly exploring and composing musical material.
Thus, in the language, you create musical "contexts" by accreting bits of harmonic and rhythmic information in a short expression,
and then to arrange those into a piece of music with combinators to play in parallel or in sequence. It also supports function-like blocks for abstracting and parameterizing repeated sequences.

In that way, composition proceeds top-down, as opposed to bottom-up approaches of other musical languages I'd encountered, where you build up note-by-note. The idea, too, was 
that you could abstract sequences and parameterize it with different harmonic information (for example: write a I-IV-V chord progression as a function, then play it in the key of C, or Db, or whichever.)

That all worked decently well. Here's the problem: Abstract ran into the fact that music is an incredibly high-dimensional space, and thus, once it was possible to sequence these little musical contexts, 
it was very hard to find any way to carve interesting musical phrases out of those blocks of marble. It was about this time that the project petered out. At best, what you have
is a little language for sequencing block chords, like an overly-basic interpretation of a jazz chart.

There are some examples in the `tunes` directory of what it looks like.

## Principles

The original principles of the project remain good. Arguably, the one we failed was #1.

1. Abstract must be able to make music you'd want to listen to.  

- If we do anything else, but not this, we've failed.

2. Abstract provides a way to compose other than laboriously "plunking down notes".

- We want to be able to write music the way we write a computer program, with modular, flexible, abstract parts, that can easily be manipulated and refactored.

3. Abstract must not be fiddly.

- We want to sit down, flip a switch, and compose.

Non-goals of the project included, above all, [livecoding](https://en.wikipedia.org/wiki/Live_coding). Abstract was supposed to be for composition, not performance.

## On chords

One thing I ran into in the course of building Abstract is that the notion of a "chord" is not very precise, and a major challenge in making a streamlined language was how to represent these. There are at least three or four kinds of things that we mean when we say "chord", all of them potentially useful in the language:

1. An "absolute" chord: a collection of concrete pitches. Cmaj: C, E, G.
2. A "diatonic" chord: in functional harmony, the triad (or beyond) that you form in the context of a scale starting from a given scale degree, possibly with modifiers. The iii or V7 chord of a particular scale.
3. A "relative" chord, i.e. a chord quality: major, minor seventh, ninth, and so on. This is a collection of intervals without reference to a particular scale.

Chords can be collections of pitches, collections of scale degrees, colletions of intervals, a chord from functional or diatonic harmony, and probably more. All of them are useful in different situations and probably ought to be represented in the language. Abstract has support for most of these, but who knows if the syntax I landed on is a good one.

This doesn't even begin to cover the concept of _voicings_ of a chord, which is another thing I don't think I found a good solution for in the language.

## Future ideas

Were I to pick this project up again, what I would focus on next is a mini-language for synthesizing musical expressions.
Again, musical expression is a really high-dimensional space, and only something like a separate language could begin to address meaningful parts of it.
That, then, could be wedded to the musical "contexts" that Abstract is currently able to sequence, and then you'd be off to the races.
(The "rhythm" expressions already in the language, which let you sequence drums to a bit pattern, maybe contain a kernel of this.)

Another thought is that it could be the wrong idea to try to address what is, in the end, a concrete collection of notes, from the analytical level down. A musical language with the kind of power Abstract was supposed to have might need a more "dialectical" approach, with the ability to move up and down the ladder from notes to harmony and harmony to notes, letting the composer specify either and then supporting them with information from the other rung on the ladder.

Rhythm is its own domain and thus probably ought to be another mini-language in the language.

The project also suffers from not being a self-contained musical environment, violating principle #3. That is, it just emits MIDI, and doesn't
generate its own audio. To hear anything, you have to do something like plug it into a JACK + MIDI environment, like running a softsynth and QJackCtl in separate windows, and wiring everything together just right. Or pipe it to outboard gear.


## Build-time dependencies 

- JACK: libjack-jackd2-dev (or libjack-dev)
- ALSA: libasound2-dev
