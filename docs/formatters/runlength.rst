Runlength
=========

Runlength is a formatter that prepends the length of the message, followed by a ":".
The actual message is formatted by a nested formatter.


Parameters
----------

**RunlengthSeparator**
  RunlengthSeparator sets the separator character placed after the runlength.
  This is set to ":" by default.

**RunlengthDataFormatter**
  RunlengthDataFormatter defines the formatter for the data transferred as message.
  By default this is set to "format.Forward" .

Example
-------

.. code-block:: yaml

	- "stream.Broadcast":
	    Formatter: "format.Runlength"
	    RunlengthSeparator: ":"
	    RunlengthFormatter: "format.Envelope"
