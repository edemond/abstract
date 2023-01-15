// Semantic analysis for Abstract. Traverses a tree of "ast" package nodes and
// emits a tree of "types" package nodes.
package main

import (
	"github.com/edemond/abstract/ast"
	"github.com/edemond/abstract/chord"
	"github.com/edemond/abstract/drivers"
	"github.com/edemond/abstract/types"
	"fmt"
	"strconv"
	"strings"
)

const ANALYZE_TRACE = false

// environment contains the name-value bindings and default part for a scope.
type environment struct {
	bindings      map[string]types.Value       // Name-value bindings for this scope.
	parameterized map[string]ast.Parameterized // Exprs we can't fully evaluate yet because of formal parameters.
	defPart       *types.SimplePart            // The default expression for this scope.
}

// Analyzer is the context object for semantic analysis of an Abstract program.
type Analyzer struct {
	// A stack of environments.
	environments []*environment
	// Instrument definitions get collected here to be
	// opened later. Key is the DeviceName.
	instruments map[string]*types.Instrument
	bpm         int
	ppq         int
	spaces      int // Spaces to indent the trace.
}

func (a *Analyzer) indent() {
	a.spaces += 1
}

func (a *Analyzer) unindent() {
	a.spaces -= 1
}

// NewAnalyzer creates a new Analyzer with default values.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		environments: make([]*environment, 0),
		instruments:  make(map[string]*types.Instrument),
		bpm:          120,
		ppq:          64,
		spaces:       0,
	}
}

// Entry point for semantic analysis.
// The root of the AST is always a PlayStatement.
func (a *Analyzer) Analyze(stmt *ast.PlayStatement) (types.Part, error) {
	a.pushScope()    // Root scope to kick things off.
	a.bindBuiltIns() // Bind built-ins in the root scope.
	return a.analyzePlay(stmt)
}

// Create Abstract's built-in scales, dynamics, and more.
func (a *Analyzer) bindBuiltIns() {
	// Scales
	major := types.NewScale([]int{2, 2, 1, 2, 2, 2, 1})
	a.bind("major", major)
	minor := types.NewScale([]int{2, 1, 2, 2, 1, 2, 2})
	a.bind("minor", minor) // TODO: Not really theoretically correct. Should be harmonic minor.
	a.bind("naturalminor", minor)
	a.bind("harmonicminor", types.NewScale([]int{2, 1, 2, 2, 1, 3, 1}))
	// TODO: Varying scales like melodic minor? Not really in scope?
	a.bind("ionian", major)
	a.bind("dorian", types.NewScale([]int{2, 1, 2, 2, 2, 1, 2}))
	a.bind("phrygian", types.NewScale([]int{1, 2, 2, 2, 1, 2, 2}))
	a.bind("lydian", types.NewScale([]int{2, 2, 2, 1, 2, 2, 1}))
	a.bind("mixolydian", types.NewScale([]int{2, 2, 1, 2, 2, 1, 2}))
	a.bind("aeolian", minor)
	a.bind("locrian", types.NewScale([]int{1, 2, 2, 1, 2, 2, 2}))
	a.bind("blues", types.NewScale([]int{3, 2, 1, 1, 3, 2}))
	a.bind("pentatonic", types.NewScale([]int{2, 2, 3, 2, 3}))

	// Dynamics (0-127)
	// TODO: hmmmmmmmmm - how to work humanization into this?
	// It's probably best to provide separate "humanize" expressions.
	a.bind("ppp", types.NewDynamics(16-1))
	a.bind("pp", types.NewDynamics(32-1))
	a.bind("p", types.NewDynamics(48-1))
	a.bind("mp", types.NewDynamics(64-1))
	a.bind("mf", types.NewDynamics(80-1))
	a.bind("f", types.NewDynamics(96-1))
	a.bind("ff", types.NewDynamics(112-1))
	a.bind("fff", types.NewDynamics(128-1))
}

// OpenInstruments opens all the instruments, using the given driver, that we found in the program text.
// Returns all the opened instruments.
func (a *Analyzer) OpenInstruments(driver drivers.Driver) ([]*types.Instrument, error) {
	a.trace("Opening instruments...")
	a.indent()
	defer a.unindent()
	insts := make([]*types.Instrument, 0)
	for name, inst := range a.instruments {
		a.trace("Opening instrument '%v'...", name)
		id, err := driver.OpenInstrument(name)
		if err != nil {
			return nil, err
		}
		inst.ID = id
		insts = append(insts, inst)
	}
	a.trace("Done opening instruments.")
	return insts, nil
}

// CloseInstruments closes all the instruments opened by the analyzer.
func (a *Analyzer) CloseInstruments(driver drivers.Driver) error {
	a.trace("Closing instruments...")
	for name, inst := range a.instruments {
		a.trace("Closing instrument '%v'...", name)
		err := driver.CloseInstrument(inst.ID)
		if err != nil {
			return err
		}
	}
	a.trace("Done closing instruments.")
	return nil
}

// trace prints a line of formatted output to stdout if tracing is enabled in the analyzer.
func (a *Analyzer) trace(s string, args ...interface{}) {
	if ANALYZE_TRACE {
		spaces := strings.Repeat(" ", a.spaces)
		fmt.Printf(spaces+s+"\n", args...)
	}
}

// Format an error with the current line number.
func (a *Analyzer) errorf(line int, err string, args ...interface{}) error {
	return fmt.Errorf("line %v: %v", line, fmt.Sprintf(err, args...))
}

// pushScope pushes a new scope onto the stack.
func (a *Analyzer) pushScope() {
	var defPart *types.SimplePart
	// If this isn't the first scope, copy the default part from the parent.
	if a.depth() > 0 {
		last := a.currentEnv()
		defPart = last.defPart.Copy()
	} else {
		defPart = types.NewSimplePart()
	}
	env := &environment{
		bindings:      make(map[string]types.Value),
		parameterized: make(map[string]ast.Parameterized),
		defPart:       defPart,
	}
	a.environments = append(a.environments, env)
}

// popScope pops the scope off the top of the stack.
func (a *Analyzer) popScope() {
	a.environments = a.environments[:a.depth()-1]
}

// Returns the depth of the current scope (how many environments we have stacked together.)
func (a *Analyzer) depth() int {
	return len(a.environments)
}

// currentEnv returns the environment at the top of the stack for checking name bindings.
func (a *Analyzer) currentEnv() *environment {
	return a.environments[a.depth()-1]
}

// contains checks if the environment has a binding for the given name, starting from
// the top-level scope and working its way down to the root scope.
func (a *Analyzer) contains(name string) bool {
	for i := a.depth() - 1; i >= 0; i-- {
		env := a.environments[i]
		_, ok := env.bindings[name]
		if ok {
			return true
		}
		_, ok = env.parameterized[name]
		if ok {
			return true
		}
	}
	return false
}

// bind binds a name to a value in the current environment.
// It panics if the name is already bound -- this condition should be checked before binding.
func (a *Analyzer) bind(name string, value types.Value) {
	env := a.currentEnv()
	_, ok := env.bindings[name]
	if ok {
		panic("Internal error: Name already bound!")
	}
	env.bindings[name] = value
}

// addParameter adds an unevaluated parameterized expression to the current scope.
func (a *Analyzer) addParameterized(name string, expr ast.Parameterized) {
	env := a.currentEnv()
	_, ok := env.parameterized[name]
	if ok {
		panic("Internal error: Parameterized expression already bound!")
	}
	env.parameterized[name] = expr
}

// analyzeLet analyzes a let statement, making sure the identifier isn't bound already,
// and creating a new binding in the environment if so.
func (a *Analyzer) analyzeLet(stmt *ast.LetStatement) error {
	if a.contains(stmt.Name) {
		return a.errorf(stmt.Line, "'%v' already defined", stmt.Name) // TODO: Where?
	}

	// So....if we have parameters in this let statement...we defer analysis until the play statement, and just
	// keep it around unanalyzed until then.
	// TODO: Later, we can optimize it so that we analyze as much as possible up front, but I don't think this will really kill us.
	p, ok := stmt.Expr.(ast.Parameterized) // TODO: Uh, this is kinda ugly, can't Expression just have the parameter checking interface?
	if ok && p.HasParameters() {
		a.trace("Storing parameterized expression '%v' for later analysis.", stmt.Name)
		a.addParameterized(stmt.Name, p)
		return nil
	}

	a.pushScope()
	value, err := a.analyzeExpr(stmt.Expr)
	if err != nil {
		a.popScope()
		return err
	}
	a.popScope()

	a.trace("Bound %v = %v", stmt.Name, value)
	a.bind(stmt.Name, value)

	return nil
}

// analyzeDefault analyzes a default statement, updating the default part in the current environment.
func (a *Analyzer) analyzeDefault(stmt *ast.DefaultStatement) error {
	a.trace("default statement")
	a.indent()
	defer a.unindent()
	simple, err := a.analyzeSimpleExpr(stmt.Expr)
	if err != nil {
		return err
	}
	part, ok := simple.(*types.SimplePart)
	if !ok {
		return a.errorf(stmt.Line, "default statement requires simple expression")
	}

	env := a.currentEnv()
	if part.Harmony.Chord.HasValue() {
		env.defPart.Harmony.Chord = part.Harmony.Chord
	}
	if part.Rhythm.Dynamics.HasValue() {
		env.defPart.Rhythm.Dynamics = part.Rhythm.Dynamics
	}
	if part.Rhythm.Humanize.HasValue() {
		env.defPart.Rhythm.Humanize = part.Rhythm.Humanize
	}
	if part.Instrument.HasValue() {
		env.defPart.Instrument = part.Instrument
	}
	if part.Interpretation.HasValue() {
		env.defPart.Interpretation = part.Interpretation
	}
	if part.Rhythm.Meter.HasValue() {
		env.defPart.Rhythm.Meter = part.Rhythm.Meter
	}
	if part.Harmony.Octave.HasValue() {
		env.defPart.Harmony.Octave = part.Harmony.Octave
	}
	if part.Harmony.Pitch.HasValue() {
		env.defPart.Harmony.Pitch = part.Harmony.Pitch
	}
	if part.Harmony.Scale.HasValue() {
		env.defPart.Harmony.Scale = part.Harmony.Scale
	}
	if part.Harmony.Voicing.HasValue() {
		env.defPart.Harmony.Voicing = part.Harmony.Voicing
	}
	return nil
}

// assertName panics if the name does not match the expected string.
func assertName(name, expected string) {
	if name != expected {
		panic(fmt.Sprintf("Internal error: argument not a %v expression", expected))
	}
}

// analyzePlay unwraps the expression contained in a "play" statement.
func (a *Analyzer) analyzePlay(stmt *ast.PlayStatement) (types.Part, error) {
	a.trace("play statement")
	a.indent()
	defer a.unindent()
	val, err := a.analyzeExpr(stmt.Expr)
	if err != nil {
		return nil, err
	}

	switch v := val.(type) {
	case *types.SimplePart:
		return v, nil
	case *types.CompoundPart:
		return v, nil
	case *types.BlockPart:
		return v, nil
	case *types.Seq: // TODO: Rename to SeqPart
		return v, nil
	case types.MessagePart:
		return v, nil
	default:
		// If an expr is just a loose value of any other kind, pack it into a SimplePart.
		simple := types.NewSimplePart()
		a.assign(simple, v)
		a.fillOutDefaults(simple)
		return simple, nil
	}
}

// analyzeOctave analyzes an octave expression (e.g. O0, O5) and returns an Octave.
func (a *Analyzer) analyzeOctave(name string) (types.Octave, error) {
	a.trace("octave expression")
	num, err := strconv.Atoi(name[1:])
	if err != nil {
		return types.NoOctave(), err
	}
	if num < 0 || num > 8 {
		return types.NoOctave(), fmt.Errorf("octave must be from 0-8 (got %v)", num)
	}
	return types.NewOctave(uint(num)), nil
}

// Analyze a pitch literal expression, like C, F#, or Gb.
func (a *Analyzer) analyzePitchLiteral(name string) (types.Pitch, error) {
	a.trace("pitch literal expression.")
	pitch, err := types.LookUpPitch(name)
	if err != nil {
		panic("Internal error: Invalid pitch, should have been caught earlier!")
	}
	return pitch, nil
}

// Analyze a rest literal expression, like _.
func (a *Analyzer) analyzeRestLiteral(name string) (types.Rest, error) {
	a.trace("rest literal expression.")
	return types.NewRest(name), nil
}

// analyzeIdentExpr looks up an identifier in the environment and returns
// its corresponding value, or an error if the identifier is not bound.
// Handles built-in identifiers like octave and pitch literals.
func (a *Analyzer) analyzeIdentExpr(expr ast.IdentExpr) (types.Value, error) {
	a.trace("identifier expression.")
	a.indent()
	defer a.unindent()
	if a.depth() <= 0 {
		panic("Internal error: Environment not set up yet!")
	}
	name := string(expr)
	if types.IsOctave(name) {
		oct, err := a.analyzeOctave(name)
		if err != nil {
			return nil, err
		}
		a.trace("'%v' evaluates to octave %v", expr, oct)
		return oct, nil
	}

	if types.IsPitch(name) {
		pitch, err := a.analyzePitchLiteral(name)
		if err != nil {
			return nil, err
		}
		a.trace("'%v' evaluates to pitch %v", expr, pitch)
		return pitch, nil
	}

	if types.IsRest(name) {
		rest, err := a.analyzeRestLiteral(name)
		if err != nil {
			return nil, err
		}
		a.trace("'%v' evaluates to rest %v", expr, rest)
		return rest, nil
	}

	// Look up the value in the environment.
	for i := a.depth() - 1; i >= 0; i-- {
		env := a.environments[i]
		val, ok := env.bindings[name]
		if ok {
			a.trace("'%v' evaluates to '%v'", expr, val)
			return val, nil
		}
		_, ok = env.parameterized[name]
		if ok {
			return nil, fmt.Errorf("missing arguments to '%v'", name) // TODO: Get line info on IdentExpr so we can report the line
		}
	}

	// It's not in the environment. Is it chord notation?
	chord, err := chord.ParseAndAnalyze(name)
	if err == nil {
		return chord, nil
	}
	fmt.Println(err)
	// TODO: Wait, what the hell to do with this error? They might not have been trying to write a chord.

	return nil, fmt.Errorf("'%v' not defined", name)
}

// analyzeMeterExpr analyzes a meter syntax sugar expression and returns a Meter.
func (a *Analyzer) analyzeMeterExpr(expr *ast.MeterExpr) (*types.Meter, error) {
	a.trace("meter expression (syntax sugar version).")
	a.indent()
	defer a.unindent()

	beats, err := a.analyzeNumberOrIdent(expr.Beats)
	if err != nil {
		return nil, err
	}
	value, err := a.analyzeNumberOrIdent(expr.Value)
	if err != nil {
		return nil, err
	}

	if int(beats.Value) <= 0 || int(value.Value) <= 0 {
		return nil, a.errorf(expr.Line, "meter: beats and value must be >= 1")
	}

	return &types.Meter{Beats: int(beats.Value), Value: int(value.Value)}, nil
}

// analyzeMeter analyzes a parameterized meter expression and returns a Meter.
func (a *Analyzer) analyzeMeter(expr *ast.ParamExpr) (*types.Meter, error) {
	a.trace("meter expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "meter")
	if len(expr.Params) != 2 {
		return nil, a.errorf(expr.Line, "meter requires meter(beats, value)")
	}

	beats, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return nil, err
	}
	value, err := a.analyzeNumberOrIdent(expr.Params[1])
	if err != nil {
		return nil, err
	}

	if int(beats.Value) <= 0 || int(value.Value) <= 0 {
		return nil, a.errorf(expr.Line, "meter: beats and value must be >= 1")
	}

	return &types.Meter{Beats: int(beats.Value), Value: int(value.Value)}, nil
}

// analyzeProb analyzes a probabalistic rhythm expression and returns a probabalistic rhythm.
func (a *Analyzer) analyzeProb(expr *ast.ParamExpr) (*types.Prob, error) {
	a.trace("prob expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "prob")
	if len(expr.Params) != 3 {
		return nil, a.errorf(expr.Line, "prob requires prob(beat, strength, percent)")
	}

	beat, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return nil, err
	}

	strength, err := a.analyzeNumberOrIdent(expr.Params[1])
	if err != nil {
		return nil, err
	}

	percent, err := a.analyzeNumberOrIdent(expr.Params[2])
	if err != nil {
		return nil, err
	}

	prob := types.NewProb(int(beat.Value), int(strength.Value), int(percent.Value))
	return prob, nil
}

// TODO: This could return a seq instead of a Velocity. Or do we not need it
// because seqs implicitly run Bjorklund?!
/*
// analyzeBjork (heh) analyzes a Bjorklund-based velocity expression and returns a ("Euclidean") Velocity.
func (a *Analyzer) analyzeBjork(expr *ast.ParamExpr) (*types.Velocity, error) {
	a.trace("Bjorklund expression.")
	assertName(expr.Name, "bjork")
	if len(expr.Params) != 2 {
		return nil, a.errorf(expr.Line, "bjork requires bjork(number pulses, number steps)")
	}

	steps, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return nil, err
	}
	pulses, err := a.analyzeNumberOrIdent(expr.Params[1])
	if err != nil {
		return nil, err
	}

	pattern := machine.Bjorklund(int(steps.Value), int(pulses.Value))
	r := types.NewVelocity()
	r.AddPartFromInts(pattern)
	return r, nil
}
*/

// analyzeScale analyzes a scale expression and returns a Scale.
func (a *Analyzer) analyzeScale(expr *ast.ParamExpr) (*types.Scale, error) {
	a.trace("scale expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "scale")
	if len(expr.Params) < 1 {
		return nil, a.errorf(expr.Line, "scale requires scale(step, ...)")
	}
	steps := make([]int, len(expr.Params))
	for i, param := range expr.Params {
		n, err := a.analyzeNumberOrIdent(param)
		if err != nil {
			return nil, err
		}
		steps[i] = int(n.Value)
	}
	return types.NewScale(steps), nil
}

// analyzeStringOrIdent analyzes a string literal or ident that must evaluate to a string.
func (a *Analyzer) analyzeStringOrIdent(expr ast.Expression) (types.String, error) {
	switch e := expr.(type) {
	case ast.StringExpr:
		return a.analyzeStringExpr(e)
	case ast.IdentExpr:
		value, err := a.analyzeIdentExpr(e)
		if err != nil {
			return "", err
		}
		s, ok := value.(types.String)
		if !ok {
			return "", fmt.Errorf("expected a string expression") // TODO: Line information
		}
		return s, nil
	default:
		return "", fmt.Errorf("expected a string")
	}
}

// analyzeNumberOrIdent analyzes a number literal or ident that must evaluate to a number.
func (a *Analyzer) analyzeNumberOrIdent(expr ast.Expression) (*types.Number, error) {
	switch e := expr.(type) {
	case *ast.NumberExpr:
		return a.analyzeNumberExpr(e)
	case ast.IdentExpr:
		value, err := a.analyzeIdentExpr(e)
		if err != nil {
			return nil, err
		}
		n, ok := value.(*types.Number)
		if !ok {
			return nil, fmt.Errorf("expected a numeric expression") // TODO: Line information
		}
		return n, nil
	default:
		return nil, fmt.Errorf("expected a number")
	}
}

// analyzeChord analyzes a relative chord expression (e.g. chord(0,4,7)) and returns a Chord.
// TODO: Account for absolute and diatonic chord expressions!
func (a *Analyzer) analyzeChord(expr *ast.ParamExpr) (types.Chord, error) {
	a.trace("chord expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "chord")
	if len(expr.Params) < 1 {
		return types.NoChord(), a.errorf(expr.Line, "chord requires at least one note")
	}
	intervals := make([]int, len(expr.Params))
	for i, param := range expr.Params {
		n, err := a.analyzeNumberOrIdent(param)
		if err != nil {
			return types.NoChord(), err
		}
		intervals[i] = int(n.Value)
	}
	return types.NewRelativeChord(1, intervals), nil // We assume chord(...) exprs are rooted on the 1 chord.
}

// analyzeDynamics analyzes a dynamics expression and returns a Dynamics.
func (a *Analyzer) analyzeDynamics(expr *ast.ParamExpr) (*types.Dynamics, error) {
	a.trace("dynamics expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "dynamics")
	ln := len(expr.Params)
	if ln != 1 && ln != 2 {
		return nil, a.errorf(expr.Line, "dynamics requires dynamics(center) or dynamics(center, humanize)")
	}
	center, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return nil, err
	}
	if int(center.Value) >= 127 || int(center.Value) < 0 {
		return nil, a.errorf(expr.Line, "center must be 0-127")
	}

	if ln == 2 {
		human, err := a.analyzeNumberOrIdent(expr.Params[1])
		if err != nil {
			return nil, err
		}
		dyn := types.NewDynamics(int(center.Value)) // TODO: Error handling!
		dyn.SetHumanize(int(human.Value))
		return dyn, nil
	} else {
		dyn := types.NewDynamics(int(center.Value)) // TODO: Error handling!
		return dyn, nil
	}
}

// analyzePitch analyzes a pitch expression and returns a Pitch.
func (a *Analyzer) analyzePitch(expr *ast.ParamExpr) (types.Pitch, error) {
	a.trace("pitch expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "pitch")
	if len(expr.Params) != 1 {
		return types.NoPitch(), a.errorf(expr.Line, "pitch requires pitch(number)")
	}
	pitch, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return types.NoPitch(), err
	}
	return types.NewPitch(pitch.Value), nil
}

// analyzePitch analyzes a voicing expression and returns a Voicing.
func (a *Analyzer) analyzeVoicing(expr *ast.ParamExpr) (types.Voicing, error) {
	a.trace("voicing expression.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "voicing")
	if len(expr.Params) != 1 {
		return types.NoVoicing(), a.errorf(expr.Line, "voicing requires voicing(number)")
	}
	voicing, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return types.NoVoicing(), err
	}
	return types.NewVoicing(voicing.Value), nil
}

// analyzeNumberExpr analyzes a numeric parameter in the context of an expression like arp(32).
func (a *Analyzer) analyzeNumberExpr(p ast.Expression) (*types.Number, error) {
	number, ok := p.(*ast.NumberExpr)
	if !ok {
		return nil, fmt.Errorf("expected numeric parameter")
	}
	return &types.Number{Value: number.Value, Digits: number.Digits}, nil
}

// analyzeStringExpr analyzes a string parameter in the context of an expression like arp(32).
func (a *Analyzer) analyzeStringExpr(p ast.Expression) (types.String, error) {
	value, ok := p.(ast.StringExpr)
	if !ok {
		return types.String(""), fmt.Errorf("expected string parameter")
	}
	return types.String(value), nil
}

// analyzeInstrument analyzes an instrument expression and returns an Instrument.
// It collects the instruments for later retrieval after the whole piece is analyzed.
func (a *Analyzer) analyzeInstrument(expr *ast.ParamExpr) (*types.Instrument, error) {
	a.trace("instrument.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "instrument")
	if len(expr.Params) != 3 {
		return nil, a.errorf(expr.Line, "instrument requires instrument(instrument name, channel, voices)")
	}
	deviceName, err := a.analyzeStringOrIdent(expr.Params[0])
	if err != nil {
		return nil, err
	}
	a.trace("Instrument name is %v", deviceName)
	channel, err := a.analyzeNumberOrIdent(expr.Params[1])
	if err != nil {
		return nil, err
	}
	if channel.Value < 1 || channel.Value > 16 {
		return nil, a.errorf(expr.Line, "instrument MIDI channel must be from 1-16")
	}
	a.trace("Instrument channel is %v", channel.Value)
	voices, err := a.analyzeNumberOrIdent(expr.Params[2])
	if err != nil {
		return nil, err
	}
	// Collect the instrument definition here.
	_, ok := a.instruments[string(deviceName)] // TODO: string conversion hack
	if ok {
		return nil, a.errorf(expr.Line, "Instrument '%v' already created", deviceName) // TODO: Where?
	}
	inst := types.NewInstrument(string(deviceName), byte(channel.Value), int(voices.Value)) // TODO: string conversion hack
	a.instruments[string(deviceName)] = inst                                                // TODO: string conversion hack
	return inst, nil
}

// analyzeNote analyzes a note expression and returns a Note.
func (a *Analyzer) analyzeNote(expr *ast.ParamExpr) (types.Note, error) {
	a.trace("note.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "note")
	if len(expr.Params) != 1 {
		return types.NoNote(), a.errorf(expr.Line, "note requires note(midi note number)")
	}
	num, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return types.NoNote(), err
	}
	note, err := types.NewNote(num.Value)
	if err != nil {
		return types.NoNote(), err
	}
	return note, nil
}

// analyzeHumanize analyzes a human(...) expression and returns a Humanize.
func (a *Analyzer) analyzeHumanize(expr *ast.ParamExpr) (*types.Humanize, error) {
	a.trace("humanize.")
	a.indent()
	defer a.unindent()
	assertName(expr.Name, "human")
	if len(expr.Params) != 1 {
		return nil, a.errorf(expr.Line, "humanize requires humanize(time)")
	}
	num, err := a.analyzeNumberOrIdent(expr.Params[0])
	if err != nil {
		return types.NoHumanize(), err
	}
	humanize, err := types.NewHumanize(num.Value)
	if err != nil {
		return types.NoHumanize(), err
	}
	return humanize, nil
}

// analyzeSeqExpr analyzes a sequence expression (like [a b c]) and returns a Seq.
func (a *Analyzer) analyzeSeqExpr(expr *ast.SeqExpr) (*types.Seq, error) {
	a.trace("sequence expression.")
	a.indent()
	defer a.unindent()
	part := types.NewSeqPart()
	part.SetParent(parent)

	parts := []types.Part{}
	for _, e := range expr.ValueExprs {

		var val types.Value
		var err error

		// TODO: Handle the (future...) built-in rest value and/or tacit value.

		switch ex := e.(type) {
		case ast.IdentExpr:
			val, err = a.analyzeIdentExpr(ex)
		case *ast.ParamExpr:
			val, err = a.analyzeParamExpr(ex)
		case *ast.MeterExpr:
			val, err = a.analyzeMeterExpr(ex)
		case *ast.SeqExpr:
			val, err = a.analyzeSeqExpr(ex)
		case ast.StringExpr:
			return nil, fmt.Errorf("loose string in sequence")
		case *ast.NumberExpr:
			return nil, fmt.Errorf("loose number in sequence")
		default:
			// Simple, compound, and block cannot appear here.
			panic("Internal error: unhandled expression type in sequence")
		}
		if err != nil {
			return nil, err
		}

		switch v := val.(type) {
		// Simple, compound, and message parts go right in.
		case *types.SimplePart:
			v.SetScale(len(expr.ValueExprs))
			parts = append(parts, v)
		case *types.CompoundPart:
			v.SetScale(len(expr.ValueExprs))
			parts = append(parts, v)
		case *types.BlockPart:
			return nil, fmt.Errorf("Sequences may not contain block parts.")
		case types.MessagePart:
			parts = append(parts, v) // TODO: Can this even happen syntactically?
		default:
			// Everything else gets upgraded to a simple part, inheriting its properties
			// first from the parent simple part in which it was found, then any other
			// properties from the default part.
			simple := types.NewSimplePart()
			a.assign(simple, v)
			a.fillOutDefaults(simple)
			simple.SetScale(len(expr.ValueExprs))
			parts = append(parts, simple)
		}
	}

	part.SetParts(parts)

	return part, nil
}

// Analyze a call to a parameterized expression.
func (a *Analyzer) analyzeCall(expr *ast.ParamExpr, p ast.Parameterized) (types.Value, error) {
	a.trace("call expression.")
	a.indent()
	defer a.unindent()

	params := p.Parameters()
	args := expr.Params

	// We should have arguments supplied for all the formal parameters.
	if len(args) != len(params) {
		return nil, a.errorf(expr.Line, "wrong number of arguments to '%v' (got %v, expected %v)", expr.Name, len(args), len(params))
	}

	// Then, analyze each argument...
	bindings := make(map[ast.IdentExpr]types.Value)
	for i := 0; i < len(params); i++ {
		val, err := a.analyzeExpr(args[i])
		if err != nil {
			return nil, err
		}
		bindings[params[i]] = val
	}

	// Finally, bind each argument to a parameter in a new environment and evaluate the expr in it.
	a.pushScope()
	defer a.popScope()
	for name, value := range bindings {
		a.trace("Binding formal parameter '%v' to '%v'.", string(name), value)
		a.bind(string(name), value)
	}
	return a.analyzeExpr(p)
}

// analyzeParamExpr analyzes/evaluates a parameterized expression, which always evaluates to a Value.
// Returns nil, err if the expression is invalid.
func (a *Analyzer) analyzeParamExpr(expr *ast.ParamExpr) (types.Value, error) {
	a.trace("parameterized expression.")
	a.indent()
	defer a.unindent()

	// First, check the environment to see if this is a function call.
	a.trace("Checking environment for parameterized expressions named '%v'.", expr.Name)
	for i := a.depth() - 1; i >= 0; i-- {
		env := a.environments[i]
		p, ok := env.parameterized[expr.Name]
		if ok {
			a.trace("'%v' parameterized expression found in the environment, analyzing now.", expr)
			return a.analyzeCall(expr, p)
		}
	}

	// Otherwise, it's a built-in.
	switch expr.Name {
	case "chord":
		return a.analyzeChord(expr)
	case "cc":
		panic("cc not implemented yet")
		//return a.analyzeCC(expr)
	case "dynamics":
		return a.analyzeDynamics(expr)
	case "human":
		return a.analyzeHumanize(expr)
	case "instrument":
		return a.analyzeInstrument(expr)
	case "meter":
		return a.analyzeMeter(expr)
	case "note":
		return a.analyzeNote(expr)
	case "pc":
		panic("pc not implemented yet")
		//return a.analyzePC(expr)
	case "pitch":
		return a.analyzePitch(expr)
	case "prob":
		return a.analyzeProb(expr)
	case "scale":
		return a.analyzeScale(expr)
	case "voicing":
		return a.analyzeVoicing(expr)
	default:
		return nil, a.errorf(expr.Line, "%v not defined", expr.Name)
	}
}

// assignAllFrom assigns all values from one SimplePart into another.
// TODO: Shouldn't this move into SimplePart itself?
func (a *Analyzer) assignAllFrom(to, from *types.SimplePart) error {
	a.trace("Assigning all values from %v to %v.", from, to)
	if to == from {
		panic("Internal error: Tried to assignAllFrom the same part to itself!")
	}

	msg := "cannot combine simple parts: %v"
	err := a.assign(to, from.Harmony.Chord)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Rhythm.Dynamics)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Rhythm.Humanize)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Instrument)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Interpretation)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Rhythm.Meter)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Harmony.Octave)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Harmony.Pitch)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Harmony.Scale)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	err = a.assign(to, from.Harmony.Voicing)
	if err != nil {
		return fmt.Errorf(msg, err)
	}
	return nil
}

// assign sets a Value on a SimplePart. If it has that Value set already, it returns an error.
// TODO: Couldn't this also move into SimplePart?
func (a *Analyzer) assign(part *types.SimplePart, value types.Value) error {
	if !value.HasValue() {
		return nil
	}
	// Switch on type of value and figure out where we can put it in the part.
	switch v := value.(type) {
	case *types.Block:
		if part.Interpretation.HasValue() {
			return fmt.Errorf("part already has interpretation %v", part.Interpretation)
		}
		part.Interpretation = v
	case types.Chord:
		if part.Harmony.Chord.HasValue() {
			return fmt.Errorf("part already has chord %v", part.Harmony.Chord)
		}
		part.Harmony.Chord = v
	case *types.Dynamics:
		if part.Rhythm.Dynamics.HasValue() {
			return fmt.Errorf("part already has dynamics %v", part.Rhythm.Dynamics)
		}
		part.Rhythm.Dynamics = v
	case *types.Humanize:
		if part.Rhythm.Humanize.HasValue() {
			return fmt.Errorf("part already has humanize %v", part.Rhythm.Humanize)
		}
		part.Rhythm.Humanize = v
	case *types.Instrument:
		if part.Instrument.HasValue() {
			return fmt.Errorf("part already has instrument %v", part.Instrument)
		}
		part.Instrument = v
	case *types.Meter:
		if part.Rhythm.Meter.HasValue() {
			return fmt.Errorf("part already has meter %v", part.Rhythm.Meter)
		}
		part.Rhythm.Meter = v
	case types.Note:
		if part.Harmony.Octave.HasValue() {
			return fmt.Errorf("part already has octave %v (tried to add note)", part.Harmony.Octave)
		}
		if part.Harmony.Pitch.HasValue() {
			return fmt.Errorf("part already has pitch %v (tried to add note)", part.Harmony.Pitch)
		}
		part.Harmony.Pitch = types.NewPitch(uint64(v % 12))
		part.Harmony.Octave = types.NewOctave(uint(v / 12))
	case types.Octave:
		if part.Harmony.Octave.HasValue() {
			return fmt.Errorf("part already has octave %v", part.Harmony.Octave)
		}
		part.Harmony.Octave = v
	case types.Pitch:
		if part.Harmony.Pitch.HasValue() {
			return fmt.Errorf("part already has pitch %v", part.Harmony.Pitch)
		}
		part.Harmony.Pitch = v
	case *types.Prob:
		if part.Interpretation.HasValue() {
			return fmt.Errorf("part already has interpretation %v", part.Interpretation)
		}
		part.Interpretation = v
	case types.Rest:
		if part.Interpretation.HasValue() {
			return fmt.Errorf("part already has interpretation %v", part.Harmony.Pitch)
		}
		part.Interpretation = v
	case *types.Scale:
		if part.Harmony.Scale.HasValue() {
			return fmt.Errorf("part already has scale %v", part.Harmony.Scale)
		}
		part.Harmony.Scale = v
	case *types.Seq:
		panic("Internal error: Can't assign a seq part to a value of a simple part.")
	case types.Voicing:
		if part.Harmony.Voicing.HasValue() {
			return fmt.Errorf("part already has voicing %v", part.Harmony.Voicing)
		}
		part.Harmony.Voicing = v
	// Unhandled stuff.
	case *types.Number:
		return fmt.Errorf("loose number in simple part")
	case types.String:
		return fmt.Errorf("loose string in simple part")
	default:
		panic(fmt.Sprintf("Internal error: unhandled value type in simple expression: %v", value))
	}
	// We trace at the bottom here because we don't want it to report
	// that we're assigning a value when in fact we haven't.
	a.trace("Assigned value %v.", value)
	return nil
}

// Take a simple part and fill out any missing values from the current environment's default expression.
func (a *Analyzer) fillOutDefaults(part *types.SimplePart) error {
	a.trace("Filling out defaults.")
	env := a.currentEnv()
	def := env.defPart
	if !part.Harmony.Chord.HasValue() {
		part.Harmony.Chord = def.Harmony.Chord
	}
	if !part.Rhythm.Dynamics.HasValue() {
		part.Rhythm.Dynamics = def.Rhythm.Dynamics
	}
	if !part.Rhythm.Humanize.HasValue() {
		part.Rhythm.Humanize = def.Rhythm.Humanize
	}
	if !part.Instrument.HasValue() {
		part.Instrument = def.Instrument
	}
	if !part.Interpretation.HasValue() {
		part.Interpretation = def.Interpretation
	}
	if !part.Rhythm.Meter.HasValue() {
		part.Rhythm.Meter = def.Rhythm.Meter
	}
	if !part.Harmony.Octave.HasValue() {
		part.Harmony.Octave = def.Harmony.Octave
	}
	if !part.Harmony.Pitch.HasValue() {
		part.Harmony.Pitch = def.Harmony.Pitch
	}
	if !part.Harmony.Scale.HasValue() {
		part.Harmony.Scale = def.Harmony.Scale
	}
	if !part.Harmony.Voicing.HasValue() {
		part.Harmony.Voicing = def.Harmony.Voicing
	}
	return nil
}

func (a *Analyzer) analyzeSimpleExpr(expr *ast.SimpleExpr) (types.Part, error) {
	a.trace("simple expression of length %v.", len(expr.ValueExprs))
	a.indent()
	defer a.unindent()
	if len(expr.ValueExprs) <= 0 {
		panic("Internal error: can't have a simple expression of length 0")
	}

	part := types.NewSimplePart()
	var seq *ast.SeqExpr

	for _, valExpr := range expr.ValueExprs {
		var val types.Value
		var err error

		// Analyze the expression...
		switch ex := valExpr.(type) {
		case ast.IdentExpr:
			val, err = a.analyzeIdentExpr(ex)
		case *ast.ParamExpr:
			// Also here, for parameterized loose values.
			// TODO: I get the sense these two cases could be deduped/simplified...
			val, err = a.analyzeParamExpr(ex)
		case *ast.MeterExpr:
			val, err = a.analyzeMeterExpr(ex)
		case *ast.SeqExpr:
			// Keep this around for later.
			// Once the rest of the simple part is analyzed, only then can we analyze the
			// seq part (because it absorbs this, the parent part's, values).
			if seq != nil {
				return nil, fmt.Errorf("more than one sequence found in simple part")
			}
			seq = ex
			continue
		default:
			panic(fmt.Sprintf("Internal error: bad ident or param, got %v", ex))
		}
		if err != nil {
			return nil, err
		}

		// Decide what to do with the resulting value.
		switch v := val.(type) {
		case *types.SimplePart:
			// TODO: NOPE. Not combining anything yet.
			a.trace("Found a simple part reference in a simple expression; combining the two.")
			err = a.assignAllFrom(part, v)
			if err != nil {
				return nil, err
			}
		case *types.CompoundPart:
			a.trace("Found a compound part reference in a simple expression.")
			if len(expr.ValueExprs) > 1 {
				fmt.Println("Warning: Expression references a compound part; these are currently ignored!")
			}
			return v, nil
		case *types.BlockPart:
			a.trace("Found a block part reference in a simple expression.")
			if len(expr.ValueExprs) > 1 {
				fmt.Println("Warning: Expression references a block part; loose values currently ignored!")
			}
			return v, nil
		case *types.Seq:
			a.trace("Found a sequence part reference in a simple expression.")
			if len(expr.ValueExprs) > 1 {
				fmt.Println("Warning: Expression references a sequence part; loose values currently ignored!")
			}
			return v, nil
		case types.MessagePart:
			// TODO: Again, I don't even think this can happen syntactically now,
			// but let's guard against it.
			a.trace("Found a message part reference in a simple expression.")
			if len(expr.ValueExprs) > 1 {
				fmt.Println("Warning: Expression references a message part; these are currently ignored!")
			}
			return v, nil
		default:
			// Here's where loose values get added to the part.
			err = a.assign(part, v)
			if err != nil {
				return nil, err
			}
		}
	}

	err := a.fillOutDefaults(part)
	if err != nil {
		return nil, err
	}

	// Okay, if at the end of this we had a seq expr in there, we can
	// analyze that now, since the parent simple part is all put together
	// to provide default values for the parts in the sequence.
	if seq != nil {
		seqPart, err := a.analyzeSeqExpr(seq, part)
		if err != nil {
			return nil, err
		}
		// A simple expr containing a seq expr produces a seq part.
		return seqPart, nil
	}

	return part, nil
}

func (a *Analyzer) analyzeCompoundExpr(expr *ast.CompoundExpr) (*types.CompoundPart, error) {
	a.trace("compound expression of %v parts", len(expr.SimpleExprs))
	a.indent()
	defer a.unindent()
	part := types.NewCompoundPart()

	// Analyze and collect each simple expression.
	for _, simple := range expr.SimpleExprs {
		p, err := a.analyzeSimpleExpr(simple)
		if err != nil {
			return nil, err
		}
		part.Add(p)
	}

	return part, nil
}

// analyzePC analyzes a PC statement and returns a part that sends a MIDI program change.
func (a *Analyzer) analyzePC(expr *ast.ParamExpr) (types.ProgramChange, error) {
	a.trace("MIDI PC expression")
	panic("yeah not yet implemented sorry")
}

func (a *Analyzer) analyzeBlockExpr(expr *ast.BlockExpr) (types.Part, error) {
	a.trace("block expression.")
	a.indent()
	defer a.unindent()
	part := types.NewBlockPart()

	// Go over its statements, building up bindings and analyzing subexpressions.
	for _, stmt := range expr.Statements {
		switch s := stmt.(type) {
		// BPM and PPQ statements can only appeare in the root scope.
		// TODO: Is how we return BPM and PPQ satisfactory? We could set it multiple times?
		case *ast.BPMStatement:
			if a.depth() > 1 {
				return nil, a.errorf(expr.Line, "bpm can only be set in the top scope")
			}
			a.trace("Setting BPM to %v.", s.BPM)
			a.bpm = s.BPM
		case *ast.PPQStatement:
			if a.depth() > 1 {
				return nil, a.errorf(expr.Line, "ppq can only be set in the top scope")
			}
			a.trace("Setting PPQ to %v.", s.PPQ)
			a.ppq = s.PPQ

			/*
			   case *ast.PCStatement:
			       if err := a.analyzePC(s); err != nil {
			           return nil, err
			       }
			   case *ast.CCStatement:
			       panic("TODO: Analyze and return a PlayStatement that sends a MIDI CC message.")
			*/

		case *ast.LetStatement:
			if err := a.analyzeLet(s); err != nil {
				return nil, fmt.Errorf("line %v: %v", s.Line, err) // TODO: uh what to do with this line business here
			}
		case *ast.DefaultStatement:
			if err := a.analyzeDefault(s); err != nil {
				return nil, err
			}
		case *ast.PlayStatement:
			p, err := a.analyzePlay(s)
			if err != nil {
				return nil, err
			}
			part.Add(p)
		default:
			panic("Internal error: unhandled statement type in block expression")
		}
	}

	a.trace("Returning block expression.")
	if part.NumParts() == 1 {
		a.trace("Block expression only contains one part, so shedding the outer block.")
		return part.FirstPart(), nil
	}
	return part, nil
}

func (a *Analyzer) analyzeExpr(expr ast.Expression) (types.Value, error) {
	// It clutters things up to put a trace here. This is just a dispatch
	// over all the types of expression, and each does their own, more
	// informative tracing.
	switch e := expr.(type) {
	case *ast.SimpleExpr:
		return a.analyzeSimpleExpr(e)
	case *ast.CompoundExpr:
		return a.analyzeCompoundExpr(e)
	case *ast.BlockExpr:
		return a.analyzeBlockExpr(e)
	case *ast.SeqExpr:
		// TODO: Verify that a seq part here, in this context, just inherits
		// the default part as its parent part.
		env := a.currentEnv()
		return a.analyzeSeqExpr(e, env.defPart)
	case ast.IdentExpr:
		return a.analyzeIdentExpr(e)
	case *ast.ParamExpr:
		return a.analyzeParamExpr(e)
	case *ast.NumberExpr:
		return a.analyzeNumberExpr(e)
	case ast.StringExpr:
		return a.analyzeStringExpr(e)
	case *ast.MeterExpr:
		return a.analyzeMeterExpr(e)
	default:
		panic("Internal error: unhandled expression type")
	}
}
