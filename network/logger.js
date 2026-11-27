
const MAX_BODY_SIZE = 1024 * 1024;

function captureResponse(r, data, flags) {
    try {
        const isSearch =
            r.method === "GET" && (
                r.uri === "/api/search" ||
                (
                    (r.uri === "/" || r.uri === "/search") &&
                    r.variables.arg_q
                )
            );

        if (!isSearch || r.variables.response_capture_error === "true") {
            return;
        }

        const body = r.rawVariables.response_body;
        const remaining = MAX_BODY_SIZE - body.length;

        if (data.length > remaining) {
            r.variables.response_truncated = "true";
        }

        if (remaining > 0) {
            r.rawVariables.response_body = Buffer.concat([
                body,
                data.subarray(0, remaining)
            ]);
        }
    } catch (error) {
        r.variables.response_capture_error = "true";
        r.error("Response capture failed: " + error.message);
    } finally {
        r.sendBuffer(data, flags);
    }
}

function responseText(r) {
    return r.rawVariables.response_body.toString("utf8");
}

export default { captureResponse, responseText };
