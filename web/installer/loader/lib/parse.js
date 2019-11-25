function parse(source, defs, verbose) {

	let lines = source.split('\n');
	lines = lines.map(function (line) {
		if (match_content(line)) {
			// console.log('match content', line);
			line = line.replace(/^\s*\/\/\/\s*#code\s*/g, '');
		}
		return line;
	});
	let hasIf = false;
	for (let n = 0; ;) {
		let startInfo = find_start_if(lines, n);
		if (startInfo === undefined) break;
		hasIf = true;
		const endLine = find_end(lines, startInfo.line);
		if (endLine === -1) {
			throw `#if without #endif in line ${startInfo.line + 1}`;
		}

		const elseLine = find_else(lines, startInfo.line, endLine);

		const cond = evaluate(startInfo.condition, startInfo.keyword, defs);

		if (cond) {
			if (verbose) {
				// console.log(`matched condition #${startInfo.keyword} ${startInfo.condition} => including lines [${startInfo.line + 1}-${endLine + 1}]`);
			}
			blank_code(lines, startInfo.line, startInfo.line);
			if (elseLine === -1) {
				blank_code(lines, endLine, endLine);
			} else {
				blank_code(lines, elseLine, endLine);
			}
		} else {
			if (elseLine === -1) {
				blank_code(lines, startInfo.line, endLine);
			} else {
				blank_code(lines, startInfo.line, elseLine);
				blank_code(lines, endLine, endLine);
			}
			if (verbose) {
				console.log(`not matched condition #${startInfo.keyword} ${startInfo.condition} => excluding lines [${startInfo.line + 1}-${endLine + 1}]`);
			}
		}

		n = startInfo.line;
	}
	let result = lines.join('\n');
	if (hasIf) {
		result = `//ifelse-loader build  ${JSON.stringify(defs)}\n` + result;
	}
	return result;
}

function match_if(line) {
	const re = /^[\s]*\/\/\/([\s]*)#(if)([\s\S]+)$/g;
	const match = re.exec(line);
	if (match) {
		return {
			line: -1,
			keyword: match[2],
			condition: match[3].trim()
		};
	}
	return undefined;
}

function match_endif(line) {
	const re = /^[\s]*\/\/\/([\s]*)#(endif)[\s]*$/g;
	const match = re.exec(line);
	return Boolean(match);
}

function match_else(line) {
	const re = /^[\s]*\/\/\/([\s]*)#(else)[\s]*$/g;
	const match = re.exec(line);
	return Boolean(match);
}

function match_content(line) {
	const re = /^[\s]*\/\/\/([\s]*)#(code)([\s\S]+)$/g;
	const match = re.exec(line);
	return Boolean(match);
}

function find_start_if(lines, n) {
	for (let t = n; t < lines.length; t++) {
		const match = match_if(lines[t]);
		if (match !== undefined) {
			match.line = t;
			return match;
			// TODO: when es7 write as: return { line: t, ...match };
		}
	}
	return undefined;
}

function find_end(lines, start) {
	let level = 1;
	for (let t = start + 1; t < lines.length; t++) {
		const mif = match_if(lines[t]);
		const mend = match_endif(lines[t]);

		if (mif) {
			level++;
		}

		if (mend) {
			level--;
			if (level === 0) {
				return t;
			}
		}
	}
	return -1;
}

function find_else(lines, start, end) {
	let level = 1;
	for (let t = start + 1; t < end; t++) {
		const mif = match_if(lines[t]);
		const melse = match_else(lines[t]);
		const mend = match_endif(lines[t]);
		if (mif) {
			level++;
		}

		if (mend) {
			level--;
		}

		if (melse && level === 1) {
			return t;
		}
	}

	return -1;
}

function blank_code(lines, start, end) {
	for (let t = start; t <= end; t++) {
		const len = lines[t].length;
		const lastChar = lines[t].charAt(len - 1);
		const windowsTermination = lastChar === '\r';
		if (len === 0) {
			lines[t] = '';
		}
		else if (len === 1) {
			lines[t] = windowsTermination ? '\r' : ' ';
		}
		else if (len === 2) {
			lines[t] = windowsTermination ? ' \r' : '//';
		}
		else {
			lines[t] = windowsTermination ? ('/').repeat(len - 1) + '\r' : ('/').repeat(len);
		}
	}
}

/**
 * @return true if block has to be preserved
 */
function evaluate(condition, keyword, defs) {

	let code = '(function(){';
	code += 'var defs = {};';
	for (let key in defs) {
		code += `defs['${key}'] = ${JSON.stringify(defs[key])};`;
	}
	code += `return (defs['${condition}']) ? true : false;})()`;

	let result;
	try {
		// console.log(code);
		result = eval(code);
		//console.log(`evaluation of (${condition}) === ${result}`);
	}
	catch (error) {
		throw `error evaluation #if condition(${condition}): ${error}`;
	}

	if (keyword === 'ifndef') {
		result = !result;
	}

	return result;
}
module.exports = parse;