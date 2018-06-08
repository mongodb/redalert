"""CLI commands for redalert"""

import click


@click.group()
def cli():
    """a system image validation tool"""
    pass


@cli.command()
def validate():
    """validate the current system"""
    click.echo("Not implemented yet!")


if __name__ == '__main__':
    cli()
