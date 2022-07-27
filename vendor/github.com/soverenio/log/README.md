# log
Contains context helpers for log

Examples:

        // initialize base context with default logger with provided trace id
        ctx, logger := log.WithTraceField(context.Background(), "TraceID")
        logger.Warn("warn")

        // get logger from context
        logger := log.FromContext(ctx)

        // initalize logger (SomeNewLogger() should return log.Logger)
        log.SetLogger(ctx, SomeNewLogger())
