import { spawn } from 'child_process';
import svelte from 'rollup-plugin-svelte';
import commonjs from '@rollup/plugin-commonjs';
import terser from '@rollup/plugin-terser';
import resolve from '@rollup/plugin-node-resolve';
import livereload from 'rollup-plugin-livereload';
import css from 'rollup-plugin-css-only';
import graphql from '@rollup/plugin-graphql'; // Use the new plugin
import json from '@rollup/plugin-json'
import replace from '@rollup/plugin-replace';

const production = !process.env.ROLLUP_WATCH;

function serve() {
	let server;

	function toExit() {
		if (server) server.kill(0);
	}

	return {
		writeBundle() {
			if (server) return;
			server = spawn('npm', ['run', 'start', '--', '--dev'], {
				stdio: ['ignore', 'inherit', 'inherit'],
				shell: true
			});

			process.on('SIGTERM', toExit);
			process.on('exit', toExit);
		}
	};
}


export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'esm', // Changed from 'iife' to 'esm'
		name: 'app',
		dir: 'public/build' // Changed from 'file' to 'dir' for ES module output
	},
	plugins: [
		svelte({
			compilerOptions: {
				// enable run-time checks when not in production
				dev: !production
			}
		}),
		// we'll extract any component CSS out into
		// a separate file - better for performance
		css({ output: 'bundle.css' }),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			dedupe: ['svelte'],
			exportConditions: ['svelte']
		}),
		commonjs(),
		graphql(), // Add the GraphQL plugin
		json(), // Add the JSON plugin to the list of plugins

		replace({
			preventAssignment: true, // This prevents Rollup from trying to rewrite imports
			'process.env.ENV': JSON.stringify(process.env.ENV),
		}),


		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload('public'),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	},
	onwarn(warning, warn) {
		// Ignore certain types of warnings or warnings from specific modules
		if (warning.code === 'CIRCULAR_DEPENDENCY' && /node_modules/.test(warning.message)) {
			return;
		}
		// For all other warnings, use the default Rollup warning handler
		warn(warning);
	}
};
