extends RichTextLabel

var timeoutDuration = 0.2
var elapsedTime = 0.0
var timeoutOccured = false
var startTime = 0
var incoming = null

func _ready() -> void:
	# create callable callback for moony_message_any signal
	var pingCallable = Callable(self, "_onMoonyMessagePing")
	# connect signal with callback
	MoonyClient.connect("moony_message_ping_ping_result", pingCallable)

func _process(delta: float) -> void:
	if not timeoutOccured:
		elapsedTime += delta
		
		if elapsedTime >= timeoutDuration:
			timeoutOccured = true
			_onTimeout()

func _onTimeout() -> void: 
	startTime = Time.get_ticks_msec()
	elapsedTime = 0
	timeoutOccured = false
	MoonyClient.sendMessage("ping", "ping", ["ping"])

func _onMoonyMessagePing(data) -> void:
	var responseTime = Time.get_ticks_msec()
	var ping = responseTime - startTime
	text = "[right]%s ms[/right]" % str(ping)
