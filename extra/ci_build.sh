#!/bin/bash

# This intended to be called from a git / svn hook (like the supplied
# examples of such).
#
# Arguments are:
#
# $1: the full path to the kerouac executable
#
# $2: the directory of the kerouac root
#
# $3: the name of this project
#
# $4: the tag of this build
#
# $5: the log file containing the description of changes represented in this
# build
#
# This script is intended, in general, to:
#
# - call kerouac with the supplied information and actually kick off the build
#
# - do any notifications / other work based on the result of that build
#
# - clean up the log file
#
# You may copy this script as a starting point to your repo and check it in,
# where it can be customized per branch (if, for example, you want
# notifications on a particular branch to be sent to a different set of
# recipients).
#

###########################################################################
# Configuration                                                           #
#                                                                         #
# (note that as written all values will pick up and prefer ENV settings). #
###########################################################################

# If MAIL_TO is set, notifications will be sent to these addresses.  If it is
# empty, no notifications will be sent.
MAIL_TO=${MAIL_TO:-""}
# MAIL_TO="person1@example.com person2@example.com"

NOTIFY_ON_FAILURE="YES"
#NOTIFY_ON_FAILURE="NO"

NOTIFY_ON_SUCCESS="YES"
#NOTIFY_ON_SUCCESS="NO"

# This command will be passed the subject and the list of MAIL_TO addresses.
MAIL_CMD=${MAIL_CMD:-"/usr/bin/env mail -s"}

# Where to find the kerouac configuration in the repo.
KEROUAC_CONFIG_NAME=${KEROUAC_CONFIG_NAME:-"kerouac.json"}

KEROUAC_BUILD_FLAGS="--remove-src"

#############
# Arguments #
#############

KEROUAC=$1
KEROUAC_ROOT=$2
PROJECT=$3
TAG=$4
LOG_FILE=$5

###############################
# Actually run the build.     #
###############################

$KEROUAC build $KEROUAC_BUILD_FLAGS . $KEROUAC_CONFIG_NAME $KEROUAC_ROOT $PROJECT $TAG

if [ $? != "0" ]
then
    STATUS=FAILED
else
    STATUS=SUCCEEDED
fi

########################################################
# Get the output from the build logs into the log file #
########################################################

echo >> $LOG_FILE
echo 'Kerouac log output:' >> $LOG_FILE
cat $($KEROUAC print kerouaclogpath $KEROUAC_ROOT $PROJECT $TAG) >> $LOG_FILE

echo >> $LOG_FILE
echo 'Build stdout:' >> $LOG_FILE
cat $($KEROUAC print stdoutpath $KEROUAC_ROOT $PROJECT $TAG) >> $LOG_FILE

echo >> $LOG_FILE
echo 'Build stderr:' >> $LOG_FILE
cat $($KEROUAC print stderrpath $KEROUAC_ROOT $PROJECT $TAG) >> $LOG_FILE

#####################################
# Send notifications if appropriate #
#####################################

if [ "$MAIL_TO" != "" ]
then
    if [ $STATUS == "FAILED" ] && [ $NOTIFY_ON_FAILURE == "YES" ]
    then
        cat $LOG_FILE | $MAIL_CMD "$PROJECT build $TAG failed" $MAIL_TO
    elif [ $NOTIFY_ON_SUCCESS == "YES" ]
    then
        cat $LOG_FILE | $MAIL_CMD "$PROJECT build $TAG succeeded" $MAIL_TO

    fi
fi

rm $LOG_FILE
