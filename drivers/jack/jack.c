#include <signal.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "_cgo_export.h"
#include "jack.h"

static jack_nframes_t frames_to_burn = 8000;

// port_buffer: 
// offset: Offset into the MIDI buffer for this period. Should be passed in verbatim from the process callback, or
// maybe later it can be modified slightly to do a "humanize" effect.
// command: Low nibble of the (command | channel) byte (e.g. 0x8 note off, 0x9 note on).
// channel: MIDI channel 1-16.
// note: MIDI note number.
// velocity: MIDI velocity.
int write_midi_event(void* port_buffer, int offset, 
    unsigned char command, unsigned char channel, 
    unsigned char note, unsigned char velocity) { 
    jack_midi_data_t* event = jack_midi_event_reserve(
        port_buffer,
        offset, // into the buffer
        3 // bytes to reserve
    );

    if (event == NULL) {
        printf("Couldn't reserve MIDI event!\n");
        return -1;
    }

    event[0] = (command << 4) | ((channel-1) & 0x0F); // note on, channel 1
    event[1] = note; // pitch
    event[2] = velocity; // velocity

    //printf("buffer: %p, event: %p, offset %d: %#x %#x %#x\n", port_buffer, event, offset, event[0], event[1], event[2]);

    return 0;
}

// Calculate the number of JACK frames for each Abstract step.
jack_nframes_t get_frames_per_step(jack_nframes_t sample_rate, int bpm, int ppq) {
    // TODO: This is not very precise. Could cause some drift, or some tempos
    // to be slightly faster or slower. But we're not planning on synchronizing
    // with anything just yet; we just want to be close for the time being.
    double bps = bpm / 60.0;
    jack_nframes_t samples_per_beat = sample_rate / bps; // TODO: Will eventually drift...
    return samples_per_beat / ppq;
}

// Main JACK callback to send some audio. We pass a jack_context* as the arg.
int process_callback(jack_nframes_t nframes, void* arg) {
    jack_context* context = (jack_context*)arg;
    jack_nframes_t start_frame = context->frames;
    PrepareBuffers(nframes);

    // This is a hack to compensate for what seems to be a problem with my JACK setup,
    // in which MIDI notes aren't actually played if I write them before a certain
    // number of frames.
    if (start_frame >= frames_to_burn) {

        // Do all of the steps that fall in this period.
        // A step happens on any multiple of context->frames_per_step that falls 
        // in the range [context->frames, (context->frames + nframes)].
        for (;;) {
            // Round from start_frame up to the next multiple of frames_per_step.
            // (y - (y % x)) + x
            jack_nframes_t next_step = 
                (start_frame - (start_frame % context->frames_per_step)) 
                + context->frames_per_step;

            if (next_step < context->frames) {
                printf("Internal error: We screwed up the math in the JACK process callback. Fix it!\n");
                return -1;
            } else if (next_step >= (context->frames + nframes)) {
                break;
            } else {
                int next_step_offset = next_step - context->frames;
                // This is the callback into Go (jack.go).
                StepSong(context->steps, next_step_offset, nframes);
                context->steps += 1;
                start_frame += context->frames_per_step;
            }
        }
    }

    context->frames += nframes;
    return 0;
}

// Do nothing. We just use this so we can continue after pause().
static void noop_signal_handler(int sig) {}

int xrun_callback() {
    printf("xrun!\n");
    return 0;
}

// Start running the driver. This will block, and the process callback 
// will be invoked, until a signal is caught.
enum jack_driver_result run_jack_driver(jack_client_t* client, int bpm, int ppq) {
    jack_nframes_t sample_rate = jack_get_sample_rate(client);

    jack_context context;
    context.frames_per_step = get_frames_per_step(sample_rate, bpm, ppq);
    context.frames = 0;
    context.steps = 0;

    int result = jack_set_process_callback(client, process_callback, &context);
    if (result != 0) {
        return JACK_SET_CALLBACK_FAILED;
    }

    result = jack_set_xrun_callback(client, xrun_callback, &context);
    if (result != 0) {
        return JACK_SET_CALLBACK_FAILED;
    }

#ifndef WIN32
    signal(SIGQUIT, noop_signal_handler);
    signal(SIGHUP, noop_signal_handler);
#endif
    signal(SIGTERM, noop_signal_handler);
    signal(SIGINT, noop_signal_handler);

    result = jack_activate(client);
    if (result != 0) {
        return JACK_ACTIVATE_FAILED;
    }

    // Sleep until we get a signal.
    pause();

    result = jack_deactivate(client);
    if (result != 0) {
        return JACK_DEACTIVATE_FAILED;
    }

    return JACK_OK;
}
