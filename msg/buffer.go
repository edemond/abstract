package msg

import (
	"fmt"
	"sort"
)

// Buffer is where Messages are collected as the song is being played.
type Buffer interface {
	Add(msg *Message)
	Any() bool
	Flip()
	Last() []*Message
	LastLength() int
	Next() []*Message
	NextLength() int
	Sort()
	Print()
}

const MAX_BUFFER_SIZE = 8192 // completely arbitrary, hopefully no one needs this many voices of polyphony

// We use double-buffering so that we know what to note-off when it's time to note-on new stuff.
// That is, at any one time step, there are two buffers for each instrument, "next" and "last".
// "next" holds the notes to play next, and "last" has the notes we just played. At each step,
// we play the notes in "next", stop the notes in "last", then just swap the buffers.
type doubleBuffer struct {
	next, last *buffer
}

// Make a note buffer per instrument, each sized to the number of voices of polyphony that instrument has.
func NewBuffer(size int) (Buffer, error) {
	if (size <= 0) || (size > MAX_BUFFER_SIZE) {
		return nil, fmt.Errorf(
			// TODO: This error message makes no sense to the user.
			"Invalid buffer size (%v). Make sure the total number of instrument voices is at least 1 and not over %v.",
			size,
			MAX_BUFFER_SIZE,
		)
	}

	return &doubleBuffer{
		next: &buffer{buf: make([]*Message, size)},
		last: &buffer{buf: make([]*Message, size)},
	}, nil
}

func (b *doubleBuffer) Add(msg *Message) {
	b.next.Add(msg)
}

func (b *doubleBuffer) Any() bool {
	return b.next.Any()
}

// TODO: this will go away when the buffer no longer has to worry about "last" vs. "next"
// orrrr....we could use slices as a "window" on the buffer, but that involves allocating a slice constantly.
func (b *doubleBuffer) LastLength() int {
	return b.last.Len()
}

func (b *doubleBuffer) Last() []*Message {
	return b.last.buf
}

// TODO: this will go away when the buffer no longer has to worry about "last" vs. "next"
// orrrr....we could use slices as a "window" on the buffer, but that involves allocating a slice constantly.
func (b *doubleBuffer) NextLength() int {
	return b.next.Len()
}

func (b *doubleBuffer) Next() []*Message {
	return b.next.buf
}

// Flip the double buffer.
func (b *doubleBuffer) Flip() {
	b.last.Clear()
	b.next, b.last = b.last, b.next
}

// Sort the buffer by HumanizeTime, increasing.
func (b *doubleBuffer) Sort() {
	// TODO: We could try sort.Stable here if we have problems with that.
	sort.Sort(b.next)
}

func (b *doubleBuffer) Print() {
	for i := 0; i < b.next.ptr; i++ {
		fmt.Printf("%v\n", b.next.buf[i])
	}
}

// buffer: one half of a doubleBuffer. Implements sort.Interface for sorting on msg.HumanizeTime.
type buffer struct {
	buf []*Message
	ptr int
}

func (b *buffer) Add(msg *Message) {
	if b.ptr < len(b.buf) {
		b.buf[b.ptr] = msg
		b.ptr += 1
	} else {
		fmt.Println("Note buffer overflow!")
	}
}

func (b *buffer) Clear() {
	b.ptr = 0
}

func (b *buffer) Any() bool {
	return b.ptr > 0
}

// Implementation of sort.Interface on buffer so that we can sort one by HumanizeTime.
// This is needed for the JACK driver, because JACK doesn't allow you to send it events
// out of order.

func (b *buffer) Len() int {
	return b.ptr
}

func (b *buffer) Less(i, j int) bool {
	// TODO: Should this take into account note off vs. note on, and ensure we're writing note off first?
	// I don't think the sort is guaranteed to be stable. We don't really have note offs here yet, though.
	return b.buf[i].HumanizeTime < b.buf[j].HumanizeTime
}

func (b *buffer) Swap(i, j int) {
	b.buf[i], b.buf[j] = b.buf[j], b.buf[i]
}
