#ifndef _ABSTRACT_JACK_H_
#define _ABSTRACT_JACK_H_
#include <jack/jack.h>
#include <jack/midiport.h>

// JACK driver context for the process callback.
typedef struct {
    jack_nframes_t frames_per_step;
    jack_nframes_t frames;
    uint64_t steps;
    jack_port_t* port;
} jack_context;

// Status codes returned from Abstract JACK driver functions where mentioned.
// If there's no other natural return value of the function we 
// return this directly, otherwise we set errno with one of these.
enum jack_driver_result {
    // If you change these, change the corresponding values in jack.go!
    JACK_OK = 0, // No error.
    JACK_SET_CALLBACK_FAILED = 1, // Failed to set the process callback.
    JACK_ACTIVATE_FAILED = 2, // Failed to activate JACK client.
    JACK_DEACTIVATE_FAILED = 3 // Failed to deactivate JACK client.
};

// Playback
extern enum jack_driver_result run_jack_driver(jack_client_t* client, int bpm, int ppq);
extern int write_midi_event(void* port_buffer, int offset, 
    unsigned char command, unsigned char channel, 
    unsigned char note, unsigned char velocity);

#endif
