#!/usr/bin/env python
# -*- coding: utf-8 -*-

from setuptools import setup, find_packages
from featness import __version__

tests_require = [
    'mock',
    'nose',
    'coverage',
    'yanc',
    'preggy',
    'tox',
    'ipdb',
    'coveralls',
    'factory_boy',
    'sphinx',
    'honcho',
]

setup(
    name='featness',
    version=__version__,
    description='Featness is a feature control application. It allows users to launch features for a selected group of users and keep track of how well those features are performing for each group.',
    long_description='''
Featness is a feature control application. It allows users to launch features for a selected group of users and keep track of how well those features are performing for each group.
''',
    keywords='lean application development web',
    author='Globo.com',
    author_email='appdev@corp.globo.com',
    url='http://featness.github.io/featness/',
    license='MIT',
    classifiers=[
        'Development Status :: 4 - Beta',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: MIT License',
        'Natural Language :: English',
        'Operating System :: MacOS',
        'Operating System :: POSIX',
        'Operating System :: Unix',
        'Operating System :: OS Independent',
        'Programming Language :: Python :: 2.7',
    ],
    packages=find_packages(),
    include_package_data=True,
    install_requires=[
        'cow-framework',
        'Flask-Admin',
        'flask-mongoengine',
        'flask-login',
        'flask-bcrypt',
        'mongoengine',
        'blinker',
    ],
    extras_require={
        'tests': tests_require,
    },
    entry_points={
        'console_scripts': [
            'featness-dashboard=featness.dashboard.main:start_server',
        ],
    },
)
