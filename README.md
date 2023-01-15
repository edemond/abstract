# Abstract

Circa 2014. This is an attempt to build a language and interpreter for music composition. The core idea was to create musical "contexts" by accreting bits of harmonic and rhythmic information in a short expression,
and then to cobble them together into a piece of music with combinators to play in parallel or in sequence. It also supports function-like blocks for abstracting and parameterizing repeated sequences.

In that way, composition proceeds top-down, as opposed to bottom-up approaches of other musical languages I'd encountered, where you build up note-by-note. The idea, too, was 
that you could abstract sequences and parameterize it with different harmonic information (for example: write a I-IV-V chord progression as a function, then play it in the key of C, or Db, or whichever.)

That all worked decently well. Here's the problem: Abstract ran into the fact that music is an incredibly high-dimensional space, and thus, once it was possible to sequence these little musical contexts, 
it was very hard to find any way to carve interesting musical phrases out of those blocks of marble. It was about this time that the project petered out. At best, what you have
is a little language for sequencing block chords, like an overly-basic interpretation of a jazz chart.


## Principles

The original principles of the project remain good. Arguably, the one we failed was #1.

1. Abstract must be able to make music you'd want to listen to.  

If we do anything else, but not this, we've failed.

2. Abstract provides a way to compose other than laboriously "plunking down notes".

We want to be able to write music the way we write a computer program, with modular, flexible, abstract parts,
that can easily be manipulated and refactored.

### 3. Abstract must not be fiddly.

We want to sit down, flip a switch, and compose.

Non-goals of the project included, above all, "livecoding". Abstract was supposed to be for composition, not performance.
	

## Future ideas

Were I to pick this project up again, what I would focus on next is a mini-language for synthesizing musical expressions.
Again, musical expression is a really high-dimensional space, and only something like a separate language could begin to address meaningful parts of it.
That, then, could be wedded to the musical "contexts" that Abstract is currently able to sequence, and then you'd be off to the races.

The project also suffers from not being a self-contained musical environment, violating principle #3. That is, it just emits MIDI, and doesn't
generate its own audio. To hear anything, you have to do something like into a fiddly JACK + MIDI environment like running a softsynth and QJackCtl in separate windows, and wiring everything together just right.


## Build-time dependencies 

- JACK: libjack-jackd2-dev (or libjack-dev)
- ALSA: libasound2-dev
