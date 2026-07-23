import * as vscode from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind,
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext) {
    const config = vscode.workspace.getConfiguration('naeos');
    const executablePath = config.get<string>('executablePath', 'naeos');

    const serverOptions: ServerOptions = {
        command: executablePath,
        args: ['lsp'],
        transport: TransportKind.stdio,
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'yaml', pattern: '**/*.neir.yaml' },
            { scheme: 'file', language: 'yaml', pattern: '**/*.naeos.yaml' },
            { scheme: 'file', language: 'yaml', pattern: '**/*.naeos.yml' },
        ],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.neir.yaml'),
        },
    };

    if (config.get<boolean>('lsp.enabled', true)) {
        client = new LanguageClient(
            'naeosLSP',
            'NAEOS NEIR Language Server',
            serverOptions,
            clientOptions,
        );
        client.start();
    }

    const validateCmd = vscode.commands.registerCommand('naeos.validate', async () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) return;
        const terminal = vscode.window.createTerminal('NEIR Validate');
        terminal.sendText(`${executablePath} validate --input "${editor.document.fileName}"`);
        terminal.show();
    });

    const suggestCmd = vscode.commands.registerCommand('naeos.suggest', async () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) return;
        const terminal = vscode.window.createTerminal('NEIR Suggest');
        terminal.sendText(`${executablePath} ai suggest --input-file "${editor.document.fileName}"`);
        terminal.show();
    });

    context.subscriptions.push(validateCmd, suggestCmd);
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) return undefined;
    return client.stop();
}
