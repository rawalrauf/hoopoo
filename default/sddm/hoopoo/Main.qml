import QtQuick 2.0
import SddmComponents 2.0

Rectangle {
  id: root
  MouseArea {
    anchors.fill: parent
    cursorShape: Qt.BlankCursor
  }
  width: 640
  height: 480
  color: "#1a1b26"

  property string currentUser: userModel.lastUser
  property bool loginFailed: false
  property bool loginInProgress: false // Keeps input elements hidden post-login
  property int sessionIndex: {
    for (var i = 0; i < sessionModel.rowCount(); i++) {
      var name = (sessionModel.data(sessionModel.index(i, 0), Qt.DisplayRole) || "").toString()
      if (name.indexOf("uwsm") !== -1)
        return i
    }
    return sessionModel.lastIndex
  }

  Connections {
    target: sddm
    function onLoginFailed() {
      root.loginFailed = true
      root.loginInProgress = false // Bring back the input fields
      password.text = ""
      password.focus = true

      // Reset the progress bar track back to 0 width for the next attempt
      barAnimation.stop()
      fillClipContainer.width = 0
    }
    function onLoginSucceeded() {
      root.loginFailed = false
    }
  }

  // --- CONFIGURABLE PROGRESS BAR TIMER ---
  Timer {
    id: loginDelayTimer
    interval: 3000 // Progress bar animation duration before login (1000 = 1s)
    repeat: false
    onTriggered: {
      sddm.login(root.currentUser, password.text, root.sessionIndex)
    }
  }

  Column {
    anchors.centerIn: parent
    spacing: 40

    Image {
      id: logo
      source: "logo.png"
      width: Math.min(sourceSize.width, root.width * 0.8)
      height: sourceSize.width > 0 ? Math.round(width * sourceSize.height / sourceSize.width) : 0
      fillMode: Image.PreserveAspectFit
      anchors.horizontalCenter: parent.horizontalCenter
    }

    Item {
      id: entryContainer
      width: entry.width
      height: entry.height
      anchors.horizontalCenter: parent.horizontalCenter

      // --- SCREEN STATE A: PASSWORD INPUT DIALOG ---
      Row {
        id: loginInputRow
        anchors.centerIn: parent
        spacing: 15
        visible: !root.loginInProgress // Stays hidden permanently after bar finishes filling

        Image {
          id: lockIcon
          source: root.loginFailed ? "lock-failed.png" : "lock.png"
          width: 34
          height: 38
          fillMode: Image.PreserveAspectFit
          anchors.verticalCenter: parent.verticalCenter
        }

        Item {
          width: entry.width
          height: entry.height

          Image {
            id: entry
            source: root.loginFailed ? "entry-failed.png" : "entry.png"
            anchors.centerIn: parent
          }

          Row {
            anchors.left: parent.left
            anchors.leftMargin: 20
            anchors.verticalCenter: parent.verticalCenter
            spacing: 5

            Repeater {
              model: Math.min(password.text.length, 21)
              Image {
                source: "bullet.png"
                width: 7
                height: 7
              }
            }
          }

          TextInput {
            id: password
            anchors.fill: parent
            anchors.leftMargin: 20
            anchors.rightMargin: 20
            verticalAlignment: TextInput.AlignVCenter
            echoMode: TextInput.Password
            font.family: "0xProto Nerd Font"
            font.pixelSize: 24
            font.letterSpacing: 5
            passwordCharacter: "\u2022"
            color: "transparent"
            selectionColor: "transparent"
            selectedTextColor: "transparent"
            cursorDelegate: Item {}
            focus: true

            onTextChanged: root.loginFailed = false

            Keys.onReturnPressed: {
              if (password.text !== "" && !root.loginInProgress) {
                root.loginInProgress = true
                barAnimation.start()
                loginDelayTimer.start()
              }
            }
            Keys.onEnterPressed: {
              if (password.text !== "" && !root.loginInProgress) {
                root.loginInProgress = true
                barAnimation.start()
                loginDelayTimer.start()
              }
            }
          }
        }
      }

      // --- SCREEN STATE B: PLYMOUTH REPLICA PROGRESS BAR ---
      Item {
        id: customProgressBarGroup
        anchors.centerIn: parent
        width: progressBoxImage.sourceSize.width
        height: progressBoxImage.sourceSize.height
        visible: root.loginInProgress

        Image {
          id: progressBoxImage
          source: "progress_box.png"
          anchors.centerIn: parent
        }

        Item {
          id: fillClipContainer
          height: progressBarImage.sourceSize.height
          width: 0
          clip: true
          anchors.left: progressBoxImage.left
          anchors.leftMargin: Math.round((progressBoxImage.sourceSize.width - progressBarImage.sourceSize.width) / 2)
          anchors.verticalCenter: progressBoxImage.verticalCenter

          Image {
            id: progressBarImage
            source: "progress_bar.png"
            anchors.left: parent.left
          }

          PropertyAnimation {
            id: barAnimation
            target: fillClipContainer
            property: "width"
            to: progressBarImage.sourceSize.width
            duration: loginDelayTimer.interval
            easing.type: Easing.Linear
          }
        }
      }
    }
  }

  Component.onCompleted: password.forceActiveFocus()
}
