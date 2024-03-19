package main

import (
	"context"
	"log"
	"strings"
	"text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/bedrock"
)

const MAX_HISTORY_LENGTH = 10
const DEFAULT_PROMPT_TEMPLATE = `
You are an assistant tasked with comparing two documents and noting the differences. These documents are manuals from medical equipment. They will be parsed with textract to provide only the raw text values. If languages other than English are present, please only evaluate and output items in English.

The individuals you are preparing the comparison for are all medical professionals and they require knowledge of the specific differences between the two so that they can effectively do their jobs.
This includes, but is not necessarily limited to, information about use, cleaning, care, and maintenance of the equipment.

Specific areas of interest include anything related to:
document ID and publish date
cleaning/sterilization
preventative maintenance
troubleshooting
testing
wiring
schematics
parts list
calibration

If the changes are only related to language, punctuation, or formatting, it is ok to simply state that, but if there are fundamental differences that a medical professional would want to know about please explicitly state them.

I will provide the documents separately. They will be inside xml tags called <document1> and <document2>. 


<document1>

3M
Ranger™
Pressure
p
A
(
Infusor
3M™Ranger
Model 145
Operators Manual
GB
Ranger™ Pressure Infusor
PT
Infusor de Pressão Ranger™,
Model 145
Modelo 145
Operators Manual
Manual do Utilizador
FR
Dispositif de perfusion sous
GR
Ranger™ Eyxuthpas Nieons
pression Ranger™ Modèle 145
MovtÉlo 145
Manuel dutilisation
XELPLOTN
DE
Ranger™ Druckinfusor Modell 145
PL
Aparat Ranger™ do wlewów
Bedienungsanleitung
cisnieniowych, model 145
IT
Infusore a pressione Ranger™
Podrecznik uzytkownika
modello 145
HU
Ranger™ infúziós pumpa,
Manuale delloperatore
145-ös modell
Felhasználói kézikönyv
ES
Infusor de presión Ranger™ modelo
145
CZ
Tlakovy infuzor Ranger™ Model 145
Manual del operador
Návod k obsluze
NL
Ranger™ -drukinfuus model 145
LT
Slèginis infuzorius "Ranger™
Bedieningshandleiding
Modelis 145
SE
Ranger™ tryckinfusor modell 145
Naudotojo vadovas
Bruksanvisning
RU
Ranger™, 145
DK
Ranger -trykinfusionsenhed
no
model 145
Betjeningsvejledning
TR
Ranger™ Basincli infüzör Model 145
Kullanim Kilavuzu
NO
Ranger™ trykkinfusjonsenhet
modell 145
CN
Ranger™ 145
Brukerhändbok
FI
RangerT-paineinfuusori, malli 145
AR
RangerTM jle.
Kayttöopas
145 jjb beially
Jewall We
3M
Ranger TM
Pressure Infusor
3M~Ranger
Model 145
(
Operators Manual
3M" Ranger
English
01
Français
19
Deutsch
39
Italiano
59
Español
79
Nederlands
99
Svenska
119
Dansk
139
Norsk
159
Suomi
179
Português
199
EAANVLKá
219
Polski
239
Magyar
259
Cesky
279
Lietuviy Kalba
299
319
Türkçe
339
359
awell well
379
iv
Revision History
RangerTM Pressure Infusor Model 145
Revision History
Revision Reason for Change
Pages Affected
Date
A
New Revision - First Release
All
May 2013
B
Translations Added
All
Jul 2013
C
Added UL 3rd edition requirements
All
Feb 2014
D
Updated to circular font and updated illustrations.
All
Sept 2016
RangerTM Pressure Infusor Model 145
Table of Contents
1
Table of Contents
Section 1: Technical Service and Order Placement
2
Technical service and order placement
2
USA
2
Proper use and maintenance
2
When you call for technical support
2
Servicing
2
Section 2: Introduction
3
Indications for use
3
Explanation of symbols
3
Explanation of signal word consequences
4
WARNING:
4
CAUTION:
5
NOTICE:
5
Product description
6
Pressure infusor panel
7
Section 3: Instructions for Use
9
Attaching the pressure infusor to the I.V. pole
9
Load and pressurize the infusors
10
Changing a fluid bag
10
Section 4: Troubleshooting
11
Standby/ON mode
11
Pressure infusor
12
Section 5: Maintenance and Storage
13
Cleaning the Ranger pressure infusor
13
Cleaning the pressure infusor and cord
13
Storage
13
Section 6: Specifications
14
Physical characteristics
16
Electrical characteristics
17
Storage and transport conditions
17
Performance characteristics
17
English 34-8719-2475-8
2
Technical Service and Order Placement
Ranger™ TM Pressure Infusor Model 145
Section 1: Technical Service and Order Placement
Technical service and order placement
USA
TEL:
1-800-228-3957
Proper use and maintenance
3M assumes no responsibility for the reliability, performance, or safety of the Ranger pressure infusor system if any of
the following events occur:
Modifications or repairs are performed by unqualified personnel.
The unit is used in a manner other than that described in the operators manual or maintenance guide.
The unit is installed in an environment that does not meet the appropriate electrical and grounding requirements.
When you call for technical support
We will need to know the serial number of your unit when you call us. The serial number label is located on the back of
the pressure infusor system.
Servicing
All service must be performed by 3M Health Care or an authorized service technician. Call 3M Health Care technical
service at 800-228-3957 for service information. Outside of the U.S. contact your local 3M Health Care representative.
Technical Service and Order Placement 34-8719-2475-8
Ranger™ Pressure Infusor Model 145
Introduction
3
Section 2: Introduction
Indications for use
The 3MTM RangerTM pressure infusor is intended to provide pressure to I.V. solution bags when rapid infusion of liquids
is required.
Explanation of symbols
The following symbols may appear on the products labeling or exterior packaging.
Power ON
Power OFF
Standby
Date of manufacture
Manufacturer
Protective earth (Ground)
CAUTION
This system is subject to the European WEEE Directive 2002/96/EC.
This product contains electrical and electronic components and must not be disposed of
using standard refuse collection. Please consult local directives for disposal of electrical and
electronic equipment.
H
Defibrillation-proof type BF applied part
VAC
Voltage, alternating current (AC)
IPX1
Drip proof protection of fluids
An equipotentiality plug (grounded) conductor other than a protective earth conductor or a neutral
conductor, providing a direct connection between electrical equipment and the potential equalization
busbar of the electrical installation. Please consult IEC 60601-1; 2005 for requirements.
CAUTION: Recycle to avoid environmental contamination
This product contains recyclable parts. For information on recycling - please contact your nearest
3M Service Center for advice.
English 34-8719-2475-8
4
Introduction
Ranger™ Pressure Infusor Model 145
Fuse
Follow instructions for use
Consult instructions for use
Keep dry
60°C
Temperature limits
-20°C
Pushing prohibited
Explanation of signal word consequences
WARNING: Indicates a hazardous situation which, if not avoided, could result in death or serious injury.
CAUTION: Indicates a hazardous situation which, if not avoided, could result in minor or
moderate injury.
NOTICE: Indicates a situation which, if not avoided, could result in property damage only.
WARNING:
1. To reduce the risks associated with hazardous voltage and fire:
Connect power cord to receptacles marked "Hospital Only," "Hospital Grade," or a reliably grounded outlet.
Do not use extension cords or multiple portable power socket outlets.
Do not position the equipment where unplugging is difficult. The plug serves as the disconnect device.
Use only the power cord specified for this product and certified for the country of use.
Examine the pressure infusor for physical damage before each use. Never operate the equipment if
the pressure infusor housing, power cord, or plug is visibly damaged. Contact 3M technical support
at 1-800-228-3957.
Do not allow the power cord to get wet.
Do not attempt to open the equipment or service the unit yourself. Contact 3M technical support at
1-800-228-3957.
Do not modify any part of the pressure infusor.
Do not pinch the power cord of the pressure infusor when attaching other devices to the I.V. pole.
2. To reduce the risks associated with exposure to biohazards:
Always perform the decontamination procedure prior to returning the pressure infusor for service and prior
to disposal.
3. To reduce the risks associated with air embolism and incorrect routing of fluids:
Never infuse fluids if air bubbles are present in the fluid line.
Ensure all luer connections are tightened.
RangerTM Pressure Infusor Model 145
Introduction
5
CAUTION:
1. To reduce the risks associated with instability, impact and facility medical device damage:
Mount the Ranger pressure infusor Model 145 on only a 3M model 90068/90124 pressure infusor I.V. pole/
base.
Do not mount this unit more than 56" (142 cm) from the floor to the base of the pressure infusor unit.
Do not use the power cord to transport or move the device.
Ensure the power cord is free from the castors during transport of the device.
Do not push on surfaces identified with the pushing prohibited symbol.
Do not pull the I.V. pole using the pressure infusor power cord.
2. To reduce the risks associated with environmental contamination:
Follow applicable regulations when disposing of this device or any of its electronic components.
3. This product is designed for pressure infusion only.
NOTICE:
1.
To avoid device damage:
Do not immerse the Ranger pressure infusor, pressure infusor parts, or accessories in any liquid or cleaning
or disinfecting solutions. The unit is not liquid proof.
Do not spray cleaning solutions onto the pressure infusor.
Do not clean the pressure infusor with abrasive cleaners or solvents such as acetone or thinner. Damage to
the case, label, and internal components may result.
Do not sterilize the pressure infusor.
2. Federal law (USA) restricts this device to sale by or on the order of a licensed healthcare professional.
3. The Ranger pressure infusor meets medical electronic interference requirements. If radio frequency interference
with other equipment should occur, connect the pressure infusor to a different power source.
4.
Device operating should include sufficient light so all labels, buttons, LEDs and interfaces can be read and/or
interpreted clearly. The operator may need to position themselves so that light shadows do not compromise the
legibility of the label on the back of the device.
English 34-8719-2475-8
6
Introduction
Ranger™ Pressure Infusor Model 145
Product description
The Ranger pressure infusor consists of the pressure infusor device, which is combined with user supplied intravenous
(IV) fluid bags which are pressurized in the infusors chambers. The Ranger pressure infusor also requires user supplied
disposable, sterile patient administration sets that can deliver the bag fluid to the patient under pressures up to 300
mmHg. The pressure infusor accepts solution bags ranging from 250 mL to 1000 mL. Fluids for use with the pressure
infusor include, but are not limited to blood, saline, sterile water, and irrigation solution. The pressure infusor is intended
to be used only with fluid bags that meet the standards of the American Association of Blood Banks. The pressure
infusor is not intended for use with fluid bags and administration sets that do not meet the above specifications.
The Ranger pressure infusor is intended for pediatric to adult patients undergoing a medical procedure. Patients may be
awake or fully anesthetized but are typically immobile on a gurney or surgical table. The fluids infused may interact with
any part of the body as determined by the medical professional. The pressure infusor system is attached to a custom I.V.
pole and the power cord is secured allowing the system to be moved around the hospital/healthcare facility to various
locations where needed. The system is designed for frequent use whenever pressure infusion is necessary.
For operation, the Ranger pressure infusor is mounted on the I.V. pole and base that is included with the pressure infusor.
A fluid bag is placed in one of the pressure infusor chambers, the infusors main power switch is turned ON, and the
infusor chambers are activated using the user interface at the front of the device. Upon activation, the air is routed to
the pressure infusor bladder, the bladder begins to inflate and pressurize the fluid bag. The user interface indicates
when the fluid bag pressure is In Range and ready for use.
Note: The solution pressure of the Ranger pressure infusor is
dependent on surface area and volume of the solution bag. To
verify the pressure, refer to the Maintenance Guide.
The Ranger pressure infusor has no user-adjustable controls.
The user slides an I.V. solution bag behind the metal fingers and
against the inflation bladder located inside the pressure infusor.
The Ranger pressure infusor should only be used in healthcare
facilities by trained medical professionals and can be used within
the patient environment.
When the infusor is attached to an external power source and the
main power switch is turned ON, pushing the pressure infusor
start/stop button ON inflates the inflation bladder and maintains
pressure on blood and solution bags.
Turning the pressure infusor start/stop button OFF deflates the
3M~Ranger
bladder. Turn the main power switch OFF when the pressure
infusor is not in use.
0
0
Main power switch
Power entry module
RangerTM Pressure Infusor Model 145
Introduction
7
Pressure infusor panel
The pressure infusor panel displays the status of the pressure infusors. When the pressure infusor is initially turned ON,
the indicators will illuminate to show operation. The pressure infusor power indicator illuminates yellow (standby) when
the main power switch is turned ON and the pressure infusors are able to be turned ON. A green LED indicates that the
infusor is ON. To pressurize/depressurize the pressure infusor, verify the pressure infusor door is closed and latched,
then press the pressure infusor start/stop button. Each pressure infusor is controlled independently.
3M TM Ranger TM
Pressure Infuser
High
High
In Range
300 mm Hg
In Range
Low
Low
1
1
1
Pressure Infusor Power
no power to unit
standby
ON
The indicator on the start/stop
A yellow indicator notifies the
A green indicator notifies the user
button notifies the user of the
user that the pressure infusor is in
that the infusor is pressurized.
status of each pressure infusor. No
standby mode and is ready to be
light indicates that the unit is either
turned ON.
not plugged in, the main power
switch is not turned ON, or there
is a system fault. See "Section 4:
Troubleshooting" on 11 for more
information.
English 34-8719-2475-8
8
Introduction
Ranger™ Pressure Infusor Model 145
3M
TM
Ranger TM
Pressure Infuser
2
High
High
3
In Range
300 mm Hg
In Range
4
Low
Low
1
2
High
Visual and audible indicator: The High yellow indicator illuminates and an audible indicator notifies
the user when the pressure infusor bladder is above 330 mmHg. The visual and audible indicator will
continue as long as the pressure remains above 330 mmHg. If the High condition is observed, the
infusor chamber should be turned OFF by using the pressure infusor start/stop button. Use of the
infusor chamber should be discontinued immediately, and 3M Patient Warming technical support
contacted for repair and servicing.
3
In Range
Visual only: The In Range green indicator flashes as the pressure is increasing in the pressure
infusor. Once the pressure is within the target range of 230-330 mmHg, the indicator will be at a
consistent green.
4
Low:
Visual and audible indicator: The Low yellow indicator illuminates and an audible indicator notifies
the user when the pressure infusor bladder has not reached 230 mmHg within approximately 30
seconds or when the pressure drops below 230 mmHg during use.
RangerTM Pressure Infusor Model 145
Instructions for Use
9
Section 3: Instructions for Use
NOTE: Assemble the model 90068/90124 pressure infusor I.V. pole/base according to the instructions for use that
accompanies the I.V. pole base.
NOTE: Assembly of the I.V. pole/base and attachment of the pressure infusor to the I.V. pole should only be performed
by a qualified, medical equipment service technician.
Attaching the pressure infusor to the I.V. pole
1. Mount the Model 145 pressure infusor onto a Model 90068/90124 pressure infusor I.V. pole/base with locking
casters as shown below.
CAUTION:
To reduce the risks associated with instability and facility medical device damage:
Do not mount this unit more than 56" (142 cm) from the floor to the base of the pressure infusor.
2. Securely fasten the clamps at the back of the infusor and tighten the knob screws until the infusor is stable.
0
3. Using the provided hook-and-loop strap, secure the power cord to the lower portion of the I.V. pole.
English 34-8719-2475-8
10
Instructions for Use
RangerTM Pressure Infusor Model 145
Load and pressurize the infusors
1. Plug the power cord into an appropriately grounded outlet.
2.
Using the main power switch located under the pressure infusor, turn the unit ON.
3. Remove excess air from I.V. fluid bags and prime.
4. Open pressure infusor door.
5. Slide fluid bag to the bottom of the pressure infusor ensuring bag is completely inside the metal fingers.
Note: Ensure solution bag port and spike hang below the pressure infusor fingers.
6. Securely close and latch the pressure infusor door.
7.
Prime the warming set. For more information about priming the set, refer to instructions provided with the
warming sets.
8. Press the pressure infusor start/stop button on the pressure infusor control panel to turn the corresponding
pressure chamber ON.
Note: A pressure infusor can only be turned ON when the start/stop button display status indicator is yellow. A green
display status indicator notifies the infusor chamber is ON.
9. After pressing the start/stop button, the infusor indicator LED should be flashing green in the In Range status.
When the indicator LED is solid green, open clamps to begin flow.
10. To depressurize, press the pressure infusor power button on the pressure infusor control panel to turn the
corresponding pressure chamber off.
Changing a fluid bag
1. Press the pressure infusor button on the pressure infusor control panel to turn the corresponding pressure
chamber OFF.
2.
Close pinch clamps on tubing.
3. Open pressure infusor door and remove fluid bag.
4. Remove the spike from the used fluid bag.
5. Remove air from the new fluid bag.
6. Insert spike into the new I.V. bag port.
7.
Push the pressure infusors bladder to expel the remaining air and slide solution bag to the bottom of the pressure
infusor, ensuring bag is completely inside the metal fingers.
Note: Ensure fluid bag port and spike hang below the pressure infusor fingers.
8.
Securely close and latch the pressure infusor door.
9.
Prime the warming set. For more information about priming the set, refer to instructions provided with the
warming sets.
10. Press the pressure infusor start/stop button on the pressure infusor control panel to turn the corresponding
pressure chamber ON.
Note: A pressure infusor can only be turned on when the start/stop button display status indicator is yellow. A green
display status indicator notifies the infusor chamber is ON.
11.
Once the pressure infusor is In Range, open clamps to resume flow from the new bag of fluid.
12. Discard fluid bags and warming sets according to institutional protocol.
RangerTM Pressure Infusor Model 145 Operators Manual
Troubleshooting
11
Section 4: Troubleshooting
All repair, calibration, and servicing of the Ranger pressure infusor requires the skill of a qualified, medical equipment
service technician who is familiar with good practice for medical unit repair. All repairs and maintenance should be in
accordance with manufacturers instructions. Service the Ranger pressure infusor every six months or whenever service
is needed. For replacement of the Ranger pressure infusor door latch, door, bladder, fingers, or power cord, contact a
biomedical technician. For additional technical support contact 3M patient warming.
Standby/ON mode
Condition
Cause
Solution
Nothing illuminates on the pressure
Power cord is not plugged into
Make sure the power cord is plugged into
infusor control panel when the
the power entry module, or
the power entry module of the pressure
main power switch is turned ON.
power cord is not plugged into
infusor. Make sure the pressure infusor is
an appropriately grounded outlet.
plugged into a properly grounded outlet.
Burned out LED light(s).
Contact a biomedical technician.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
Power LED status lights do
Power cord is not plugged into
Make sure the power cord is plugged into
not illuminate.
the power entry module, or
the power entry module of the pressure
power cord is not plugged into
infusor. Make sure the pressure infusor is
an appropriately grounded
plugged into a properly grounded outlet.
outlet.
Unit is not turned ON.
Using the main power switch located under
the pressure infusor turn the unit ON.
Burned out LED light.
Push the pressure infusor start/stop button,
if unit functions properly, continue use.
Contact a biomedical technician after use to
replace LED.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
The status indicators (Low, In
Power cord is not plugged into
Make sure the power cord is plugged into
Range, and/or High) do not
the power entry module, or
the power entry module of the pressure
illuminate when the pressure
power cord is not plugged into
infusor. Make sure the pressure infusor is
infusor start/stop button is pushed.
an appropriately grounded outlet.
plugged into a properly grounded outlet.
Unit is not turned ON.
Using the main power switch located under
the pressure infusor turn the unit ON.
Burned out LED light(s).
Contact a biomedical technician.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
Troubleshooting 34-8719-2475-8
12
Troubleshooting
RangerTM Pressure Infusor Model 145
Pressure infusor
Condition
Cause
Solution
Pressure infusor is not
Power cord is not plugged into the power
Make sure the power cord is plugged into the
working.
entry module, or power cord is not
power entry module of the pressure infusor.
plugged into an appropriately grounded
Make sure the pressure infusor is plugged into
outlet.
a properly grounded outlet.
Unit is not turned ON.
Using the main power switch located under
the pressure infusor turn the unit ON.
Unit fault.
Discontinue use of unit. Contact a
biomedical technician.
Blown fuse.
Contact a biomedical technician.
Low indicator
Pressure infusor bladder is loose or has
Discontinue use of pressure infusor
(solid yellow visual with
become unattached.
chamber. Use the other side of the
audible indicator).
pressure infusor
Reattach bladder by using your thumbs
to fit one side of the bladder port on the
bladder retaining collar and stretch into
position.
Pressure infusor door may not be closed
Securely close and latch the pressure infusor
and securely latched.
door.
Detected pressure has fallen below 230
Continue infusion or use the other side of
mmHg.
the pressure infusor. Contact a biomedical
technician after use.
High indicator
Pressure is above 330 mmHg.
Discontinue use of pressure infusor chamber.
(solid yellow visual with
Use the other side of the pressure infusor.
audible indicator).
Contact a biomedical technician after use.
Leakage of fluid.
Bag is not spiked securely.
Secure spike in bag.
Bladder does not
Unit fault.
Contact a biomedical technician after use.
deflate after pressure is
discontinued.
RangerTM Pressure Infusor Model 145
Maintenance and Storage
13
Section 5: Maintenance and Storage
Cleaning the Ranger pressure infusor
Clean the Ranger pressure infusor on an as-needed basis or per institutional protocol.
NOTICE:
Do not immerse the pressure infusor or use a dripping wet cloth for cleaning. Moisture may seep inside the
device and damage the electrical components. Clean the pressure infusor exterior with a soft cloth using
plain water or a mild, all-purpose or nonabrasive cleaner.
Do not spray cleaning solutions onto the pressure infusor. Damage to the pressure infusor may occur.
Cleaning the pressure infusor and cord
Clean the Ranger pressure infusor on an as-needed basis or per facility policies and procedures for cleaning
electronic equipment.
1. Disconnect power supply from power outlet.
2.
Use a slightly damp soft cloth moistened with a mild, nonabrasive cleaning solution to clean the device surfaces,
hose, and cord. Avoid getting liquid into electronic ports.
3. Dry with a separate soft cloth.
Storage
Cover and store all components in a cool, dry place when not in use. Take care not to drop or jar the unit.
English 34-8719-2475-8
14
Specifications
RangerTM Pressure Infusor Model 145
Section 6: Specifications
GUIDANCE AND MANUFACTURERS DECLARATION - ELECTROMAGNETIC EMISSIONS
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the
model 145 should assure that it is used in such an environment.
EMISSIONS TEST
COMPLIANCE
ELECTROMAGNETIC ENVIRONMENT - GUIDANCE
RF emissions
Group 1
The model 145 uses RF energy only for its internal function.
CISPR 11
Therefore, its RF emissions are very low and are not likely to
cause any interference in nearby electronic equipment.
RF emissions
Class B
The model 145 is suitable for use in all establishments, including
CISPR 11
domestic establishments and those directly connected to
the public low-voltage power supply network that supplies
Harmonic emissions
Class A
buildings used for domestic purposes.
IEC 61000-3-2
Voltage fluctuations/Flicker
Complies
emissions
IEC 61000-3-3
GUIDANCE AND MANUFACTURERS DECLARATION - ELECTROMAGNETIC IMMUNITY
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the
model 145 should assure that it is used in such an environment.
IMMUNITY TEST
IEC 60601
COMPLIANCE
ELECTROMAGNETIC ENVIRONMENT - GUIDANCE
TEST LEVEL
LEVEL
Electrostatic discharge
=6 kV contact
+6 kV contact
Floors should be wood, concrete or ceramic tile.
(ESD)
+8 kV air
+8 kV air
If floors are covered with synthetic material, the
IEC 61000-4-2
relative humidity should be at least 30%.
Electrical fast transient/
2 kV power lines
2 kV power
Mains power quality should be that of a typical
burst
lines
commercial or hospital environment.
IEC 61000-4-4
Surge
1 kV line to line
1 kV line to line
Mains power quality should be that of a typical
IEC 61000-4-5
2 kV line to gnd
2 kV line to
commercial or hospital environment.
gnd
Voltage dips, short
<5% UT
<5% UT
Mains power quality should be that of a typical
interruptions and voltage
(>95% dip in UT)
(>95% dip in UT)
commercial or hospital environment. If the user of
variations on power
for 0,5 cycle
for 0,5 cycle
the model 145 requires continued operation during
supply input lines
40% UT
40% UT
power mains interruptions, it is recommended that
IEC 61000-4-11
(60% dip in UT)
(60% dip in UT)
the model 145 be powered from an uninterruptible
for 6 cycles
for 6 cycles
power supply or a battery.
70% UT
70% U-
(30% dip in UT)
(30% dip in UT)
for 30 cycles
for 30 cycles
<5% UT
<5% UT
(>95% dip in UT)
(>95% dip in UT)
for 5 sec
for 5 sec
Power frequency (50/60
3 A/m
3 A/m
Power frequency magnetic fields should be at
Hz) magnetic field
levels characteristic of a typical location in a typical
IEC 61000-4-8
commercial or hospital environment.
NOTE UT is the a.c. mains voltage prior to application of the test level.
GUIDANCE AND MANUFACTURERS DECLARATION - ELECTROMAGNETIC IMMUNITY
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the
model 145 should assure that it is used in such an environment.
IMMUNITY TEST
IEC 60601
COMPLIANCE
ELECTROMAGNETIC ENVIRONMENT - GUIDANCE
TEST LEVEL
LEVEL
Portable and mobile RF communications equipment should be
used no closer to any part of the model 145, including cables,
than the recommended separation distance calculated from
the equation applicable to the frequency of the transmitter.
Recommended separation distance
d = 1.2 VP
d = 1.2 VP 80 MHz to 800 MHz
Conducted RF
3 Vrms
d = 2.3 VP 800 MHz to 2,5 GHz
IEC 61000-4-6
150 kHz to 80 MHz
3 Vrms
where P is the maximum output power rating of the
transmitter in watts (W) according to the transmitter
Radiated RF
3 V/m
3 V/m
manufacturer and d is the recommended separation distance
IEC 61000-4-3
80 MHz to 2,5 GHz
in meters (m).
Field strengths from fixed RF transmitters, as determined
by an electromagnetic site survey, should be less than the
compliance level in each frequency range.b
Interference may occur in the vicinity of equipment marked
with the following symbol:
((())))
NOTE 1 At 80 MHz and 800 MHz, the higher frequency range applies.
NOTE 2 These guidelines may not apply in all situations. Electromagnetic propagation is affected by absorption and
reflection from structures, objects and people.
Field strengths from fixed transmitters, such as base stations for radio (cellular/cordless) telephones and land
mobile radios, amateur radio, AM and FM radio broadcast and TV broadcast cannot be predicted theoretically with
accuracy. To assess the electromagnetic environment due to fixed RF transmitters, an electromagnetic site survey
should be considered. If the measured field strength in the location in which the model 145 is used exceeds the
applicable RF compliance level above, the model 145 should be observed to verify normal operation. If abnormal
performance is observed, additional measures may be necessary, such as re-orienting or relocating the model 145.
b Over the frequency range 150 kHz to 80 MHz, field strengths should be less than 3 V/m.
16
Specifications
RangerTM Pressure Infusor Model 145
RECOMMENDED SEPARATION DISTANCES BETWEEN PORTABLE AND MOBILE RF COMMUNICATIONS EQUIPMENT AND THE MODEL 145
The model 145 is intended for use in an electromagnetic environment in which radiated RF disturbances are
controlled. The customer or the user of the model 145 can help prevent electromagnetic interference by maintaining
a
minimum distance between portable and mobile RF communications equipment (transmitters) and the model 145 as
recommended below, according to the maximum output power of the communications equipment.
SEPARATION DISTANCE ACCORDING TO FREQUENCY OF TRANSMITTER
RATED MAXIMUM OUTPUT
m
POWER OF TRANSMITTER
W
150 kHz TO 80 MHz
80 MHz TO 800 MHz
800 MHz TO 2,5 GHz
d = 1.2 VP
d = 1.2 VP
d = 2.3 VP
0,01
0.12
0.12
0.23
0,1
0.37
0.37
0.74
1
1.17
1.17
2.30
10
3.69
3.69
7.37
100
11.67
11.67
23.30
For transmitters rated at a maximum output power not listed above, the recommended separation distance d in meters
(m) can be estimated using the equation applicable to the frequency of the transmitter, where P is the maximum output
power rating of the transmitter in watts (W) according to the transmitter manufacturer.
NOTE 1 At 80 MHz and 800 MHz, the separation distance for the higher frequency range applies.
NOTE 2 These guidelines may not apply in all situations. Electromagnetic propagation is affected by absorption and
reflection from structures, objects and people.
Physical characteristics
Dimensions
15.75 in (40 cm) high
20 in (51 cm) wide, 7.75 in (20 cm) deep
Weight
17 lb. (7.7kg)
Mounting
Dual clamp
Classification
Classified under IEC 60601-1 Guidelines as Class 1, Type BF, Ordinary Equipment.
Medical Equipment 4HZ8
CLASSIFIED
MEDICAL - GENERAL MEDICAL EQUIPMENT AS TO ELECTRICAL SHOCK, FIRE AND MECHANICAL
HAZARDS ONLY IN ACCORDANCE WITH UL 60601-1; CAN/CSA-C22.2, No.601.1; ANSI/AAMI
C
US
ES60601-1:2005 CSA-C22.2 No. 60601-1:08; Control No.4HZ8
Classified under IEC 60601-1 Guidelines (and other national versions of the Guidelines) as Class I, Type BF, defibrillation
proof, IPX1 rated, Ordinary equipment, Continuous operation. Classified by Underwriters Laboratories Inc. with respect
to electric shock, fire and mechanical hazards only, in accordance with IEC/EN 60601-1 and in accordance with
Canadian/CSA C22.2, No. 601.1. Classified under the Medical Device Directive as a Class llb device.
RangerTM Pressure Infusor Model 145
Specifications
17
Electrical characteristics
Leakage current
Meets leakage current requirements in accordance with AAMI 60601-1 and IEC 60601-1.
Power cord
15 feet (4.6 m)
Device rating
110-120 VAC, 50/60 Hz, 1 Amp
220-240 VAC, 50/60 Hz, 0.8 Amp
Fuse
2 X F1A-H, rated 250V, for 110-120 VAC unit
2 X F0.8A-H, rated 250V, for 220-240 VAC unit
Storage and transport conditions
Storage/transport temperature
-20 to 60°C (-4°F to 140°F)
Store all components at room temperature and in a dry place when in use.
Operating humidity
Up to 90% RH, noncondensing
Atmospheric pressure range
Altitude up to 2000m or 80 kPa
Performance characteristics
Operating pressure
300 mmHg setpoint
Note:
Pressure system is In Range when the pressure infusor bladders are inflated to between 230 mmHg (low)
and 330 mmHg (high). If pressure falls below 230 mmHg for more than approximately 30 seconds the Low
yellow indicator will illuminate and an audible indicator will sound. The High yellow and audible indicator
notifies the user when the pressure infusor bladder is above 330 mmHg.
The outlet pressure of the fluid may vary with the surface area and volume of the fluid bag.
English 34-8719-2475-8
3M
Made in the USA of globally
sourced material by
wl 3M Health Care
EC|REP
3M Deutschland GmbH
2510 Conway Ave.
Health Care Business
St. Paul, MN 55144 U.S.A.
Carl-Schurz-Str. 1
1-800-228-3957
41453 Neuss
www.rangerfluidwarming.com
Germany
3M and RANGER are trademarks of 3M.
Used under license in Canada.
© 2016, 3M. All rights reserved.
Issue Date: 2016-09
34-8719-2475-8
</document1>


<document2>

3M
Ranger
C
N
3M™ Ranger
en
Pressure Infusor Model 145 Operators Manual
02
fr
Manuel dutilisation du dispositif de perfusion sous
pression modèle 145
12
de
Druckinfusor, Modell 145 - Bedienungsanleitung
24
it
Manuale per loperatore dellinfusore a pressione modello 145
35
es
Manual del operador del infusor de presión modelo 145
46
nl
Pressure Infusor Model 145 Bedieningshandleiding
57
SV
Tryckinfusor modell 145 användarhandbok
68
da
Brugsanvisning til trykinfusionsenhed model 145
79
no
Trykkinfusjonsenhet modell 145, Brukerhändbok
90
fi
Paineinfuusorin malli 145 käyttöopas
101
pt
Manual do Operador do Infusor de Pressão, Modelo 145
112
el
Eyespiolo Eyxutnpa Nisons Movtálo 145
123
pl
Infuzor cisnieniowy, model 145 Instrukcja obstugi
135
hu
Infúziós pumpa, 145-ös típus - Kezelöi kézikönyv
147
CS
Tlakovy infuzor, model 145 - návod k obsluze
158
sk
Príruka pre obsluhu pre tlakovy infúzor, model 145
169
sl
Navodila za uporabo tlacnega infuzorja, model 145
180
et
Röhkinfuusori mudeli 145 kasutusjuhend
191
Iv
Spiediena infüzijas ierice, modelis 145,
lietoanas rokasgrämata
202
It
Slègio infuzoriaus modelio 145 naudotojo vadovas
213
ro
Perfuzor sub presiune Model 145 Manual de utilizare
224
ru
KOMNo
, 145
236
hr
Korisniki prirunik za infuzijsku pumpu, model 145
249
bg
3a
260
tr
Basinçle infüzör Model 145 Kullanim Kilavuzu
272
zh
145
AII
283
ar
bisle
145 jibb Series JWS
293
Table of Contents
en
Section 1: Technical Service and Order Placement
3
Technical service and order placement
3
USA
3
Proper use and maintenance
3
When you call for technical support
3
Servicing
3
Section 2: Introduction
3
Product description
3
Indications for use
3
Patient Population and Settings
3
Explanation of signal word consequences
3
WARNING:
3
CAUTION:
4
NOTICE:
4
Overview and Operation
4
Pressure infusor panel
5
Section 3: Instructions for Use
6
Attaching the pressure infusor to the I.V. pole
6
Load and pressurize the infusors
6
Changing a fluid bag
6
Section 4: Troubleshooting
7
Standby/ON mode
7
Pressure infusor
7
Section 5: Maintenance and Storage
8
General Maintenance and Storage
8
Cleaning Instructions
8
Storage
8
Servicing
8
Symbol Glossary
8
Section 6: Specifications
9
Physical characteristics
10
Electrical characteristics
11
Storage and transport conditions
11
Performance characteristics
11
2
Section 1: Technical Service and Order Placement
Technical service and order placement
USA: TEL: 1-800-228-3957 (USA Only)
Outside of the USA: Contact your local 3M representative.
Proper use and maintenance
3M assumes no responsibility for the reliability, performance, or safety of the Ranger pressure infusor system if any of the following events occur:
Modifications or repairs are not performed by a qualified, medical equipment service technician who is familiar with good practice for medical device repair.
The unit is used in a manner other than that described in the Operators or Preventive Maintenance Manual.
The unit is installed in an environment that does not provide grounded electrical outlets.
The warming unit is not maintained in accordance with the procedures described in the Preventive Maintenance Manual.
When you call for technical support
We will need to know the serial number of your unit when you call us. The serial number label is located on the back of the pressure infusor system.
Servicing
All service must be performed by 3M or an authorized service technician. Call 3M at 1-800 228-3957 (USA only) for service information. Outside of the USA,
contact your local 3M representative.
Section 2: Introduction
Product description
The Ranger pressure infusor consists of the pressure infusor device, which is combined with user supplied intravenous (IV) fluid bags which are pressurized in the
infusors chambers. The Ranger pressure infusor also requires user supplied disposable, sterile patient administration sets that can deliver the bag fluid to the patient
under pressures up to 300 mmHg. The pressure infusor accepts solution bags ranging from 250 mL to 1000 mL. Fluids for use with the pressure infusor include, but
are not limited to blood, saline, sterile water, and irrigation solution. The pressure infusor is intended to be used only with fluid bags that meet the standards of the
American Association of Blood Banks. The pressure infusor is not intended for use with fluid bags and administration sets that do not meet the above specifications.
Indications for use
The 3MTM RangerTM pressure infusor is intended to provide pressure to I.V. solution bags when rapid infusion of liquids is required.
Patient Population and Settings
Adult and pediatric patients being treated in operating rooms, emergency trauma settings, or other areas when rapid infusion of liquids is required. The
liquids infused may interact with any part of the body as determined by the medical professional.
Explanation of signal word consequences
WARNING: Indicates a hazardous situation which, if not avoided, could result in death or serious injury.
CAUTION: Indicates a hazardous situation which, if not avoided, could result in minor or moderate injury.
NOTICE: Indicates a situation which, if not avoided, could result in property damage only.
WARNING:
1. To reduce the risks associated with hazardous voltage and fire:
Connect power cord to receptacles marked "Hospital Only," "Hospital Grade," or a reliably grounded outlet. Do not use extension cords or multiple
portable power socket outlets.
Do not position the equipment where unplugging is difficult. The plug serves as the disconnect device.
Use only the power cord specified for this product and certified for the country of use.
Examine the pressure infusor for physical damage before each use. Never operate the equipment if the pressure infusor housing, power cord, or plug
is visibly damaged. In the USA, contact 3M at 1-800-228-3957 (USA only). Outside of the USA contact your local 3M representative.
Do not allow the power cord to get wet.
Do not attempt to open the equipment or service the unit yourself. In the USA, contact 3M at 1-800-228-3957 (USA only). Outside of the USA
contact your local 3M representative.
Do not modify any part of the pressure infusor.
Do not pinch the power cord of the pressure infusor when attaching other devices to the I.V. pole.
2. To reduce the risks associated with exposure to biohazards:
Always perform the decontamination procedure prior to returning the pressure infusor for service and prior to disposal.
3. To reduce the risks associated with air embolism and incorrect routing of fluids:
Never infuse fluids if air bubbles are present in the fluid line.
Ensure all luer connections are tightened.
3
CAUTION:
1. To reduce the risks associated with instability, impact and facility medical device damage:
Mount the Ranger pressure infusor Model 145 only on a 3M model 90068/90124 pressure infusor I.V. pole/base.
Do not mount this unit more than 56" (142 cm) from the floor to the base of the pressure infusor unit.
Do not use the power cord to transport or move the device.
Ensure the power cord is free from the castors during transport of the device.
Do not push on surfaces identified with the pushing prohibited symbol.
Do not pull the I.V. pole using the pressure infusor power cord.
2. To reduce the risks associated with environmental contamination:
Follow applicable regulations when disposing of this device or any of its electronic components.
3. This product is designed for pressure infusion only.
NOTICE:
1. The Ranger pressure infusor meets medical electronic interference requirements. If radio frequency interference with other equipment should occur,
connect the pressure infusor to a different power source.
2. Do not use cleaning solutions with greater than 80% alcohol or solvents, including acetone and thinner, to clean the warming unit or hose. Solvents may
damage the labels and other plastic parts.
3. Do not immerse the Ranger unit or accessories in any liquid or subject them to any sterilization process.
Overview and Operation
Overview and Operation
The pressure infusor system is attached to a custom I.V. pole with base. The power cord is secured allowing the system to be moved around the hospital/healthcare
facility to various locations where needed. The system is designed for frequent use whenever pressure infusion is necessary.
For operation, the Ranger pressure infusor is mounted on the custom I.V. pole and base. A fluid bag is placed in one of the pressure infusor chambers, the infusors main
power switch is turned ON, and the infusor chambers are activated using the user interface at the front of the device. Upon activation, the air is routed to the pressure
infusor bladder, the bladder begins to inflate and pressurize the fluid bag. The user interface indicates when the fluid bag pressure is In Range and ready for use.
Note: The solution pressure of the Ranger pressure infusor is dependent on surface area and volume of the solution bag. To verify the pressure, refer to the preventive
maintenance manual.
The Ranger pressure infusor has no user adjustable controls. The user slides an I.V. solution bag behind the metal fingers and against the inflation bladder located
inside the pressure infusor. The Ranger pressure infusor should only be used in healthcare facilities by trained medical professionals and can be used within the patient
environment.
When the infusor is attached to an external power source and the main power switch is turned ON, pushing the pressure infusor start/stop button ON inflates the
inflation bladder and maintains pressure on blood and solution bags.
Turning the pressure infusor start/stop button OFF deflates the bladder. Turn the main power switch OFF when the pressure infusor is not in use
0
0
Main power switch
3M" Ranger™
Power entry module
Power entry module
Main power switch
4
Pressure infusor panel
The pressure infusor panel displays the status of the pressure infusors.
When the pressure infusor is initially turned ON, the indicators will
illuminate to show operation. The pressure infusor power indicator
3M
Ranger™
illuminates yellow (standby) when the main power switch is turned
ON and the pressure infusors are able to be turned ON. A green
LED indicates that the infusor is ON. To pressurize/depressurize
the pressure infusor, verify the pressure infusor door is closed and
High
High
latched, then press the pressure infusor start/stop button. Each
pressure infusor is controlled independently.
In Range
300 mmHg
In Range
Low
Low
1
1
Pressure Infusor Power
no power to unit
standby
ON
The indicator on the start/stop button notifies the user of the status
A yellow indicator notifies the user
A green indicator notifies the user that
of each pressure infusor. No light indicates that the unit is either
that the pressure infusor is in standby
the infusor is pressurized.
not plugged in, the main power switch is not turned ON, or there
mode and is ready to be turned ON.
is a system fault. See "Section 4: Troubleshooting" on 7 for more
information.
2
3M
Ranger™
TM
High
Visual and audible indicator: The High yellow indicator illuminates and
an audible indicator notifies the user when the pressure infusor bladder is
2
above 330 mmHg. The visual and audible indicator will continue as long as
the pressure remains above 330 mmHg. If the High condition is observed,
High
High
the infusor chamber should be turned OFF by using the pressure infusor
start/stop button. Use of the infusor chamber should be discontinued
3
immediately, and 3M contacted for repair and servicing.
In Range
300 mmHg
In Range
3
4
In Range
Low
Low
Visual only: The In Range green indicator flashes as the pressure is
increasing in the pressure infusor. Once the pressure is within the target
range of 230 330 mmHg, the indicator will be at a consistent green.
4
Low
Visual and audible indicator: The Low yellow indicator illuminates and an
audible indicator notifies the user when the pressure infusor bladder has not
reached 230 mmHg within approximately 30 seconds or when the pressure
drops below 230 mmHg during use.
5
Section 3: Instructions for Use
NOTE: Assemble the model 90068/90124 pressure infusor I.V. pole/base according to the instructions for use that accompanies the I.V. pole base.
NOTE: Assembly of the I.V. pole/base and attachment of the pressure infusor to the I.V. pole should only be performed by a qualified, medical equipment
service technician.
Attaching the pressure infusor to the I.V. pole
1. Mount the Model 145 pressure infusor onto a Model 90068/90124 pressure
infusor custom I.V. pole and base.
0
CAUTION:
To reduce the risks associated with instability and facility medical device damage:
Do not mount this unit more than 56" (142 cm) from the floor to the base
of the pressure infusor.
2. Securely fasten the clamps at the back of the infusor and tighten the knob
screws until the infusor is stable.
3. Using the provided hook and loop strap, secure the power cord to the lower
portion of the I.V. pole.
Load and pressurize the infusors
1. Plug the power cord into an appropriately grounded outlet.
2. Using the main power switch located under the pressure infusor, turn the unit ON.
3. If administration set is used, spike bags and prime administration set tubing, ensuring that all air is removed from the tubing set.
4. Slide the cassette into the slot in the 3MTM RangerTM Warming Unit, Model 245 or 247. The cassette can only fit into the device one way
5.
Connect warming set and continue priming, ensuring that all air is removed from the warming set. If no administration set is used, spike bags and prime warming
set ensuring that all air is removed from the tubing set. For more information about priming the set, refer to instructions provided with the warming sets.
6. Open pressure infusor door.
7. Slide fluid bag to the bottom of the pressure infusor ensuring bag is completely inside the metal fingers.
Note: Ensure solution bag port and spike hang below the pressure infusor fingers.
8. Securely close and latch the pressure infusor door.
9. Close pinch clamps on tubing
10. Press the pressure infusor start/stop button on the pressure infusor control panel to turn the corresponding pressure chamber ON.
Note: A pressure infusor can only be turned ON when the start/stop button display status indicator is yellow. A green display status indicator notifies the
infusor chamber is ON.
11. After pressing the start/stop button, the infusor indicator LED should be flashing green in the In Range status. When the indicator LED is solid green,
open clamps to begin flow.
12. To depressurize, press the pressure infusor power button on the pressure infusor control panel to turn the corresponding pressure chamber off.
Changing a fluid bag
1.
Press the pressure infusor button on the pressure infusor control panel to turn the corresponding pressure chamber OFF.
2. Close pinch clamps on tubing.
3. Open pressure infusor door and remove fluid bag.
4. Remove the spike from the used fluid bag.
5. Insert spike into the new I.V. bag port.
6.
Push the pressure infusors bladder to expel the remaining air and slide solution bag to the bottom of the pressure infusor, ensuring bag is completely
inside the metal fingers.
Note: Ensure fluid bag port and spike hang below the pressure infusor fingers.
7.
Securely close and latch the pressure infusor door.
8. Prime the warming set, ensuring all air is removed from the tubing. For more information about priming the set, refer to instructions provided with the
warming sets.
9. Press the pressure infusor start/stop button on the pressure infusor control panel to turn the corresponding pressure chamber ON.
Note: A pressure infusor can only be turned on when the start/stop button display status indicator is yellow. A green display status indicator notifies the
infusor chamber is ON.
10. Once the pressure infusor is In Range, open clamps to resume flow from the new bag of fluid.
11. Discard fluid bags and warming sets according to institutional protocol.
6
Section 4: Troubleshooting
All repair, calibration, and servicing of the Ranger pressure infusor requires the skill of a qualified, medical equipment service technician who is familiar with
good practice for medical unit repair. All repairs and maintenance should be in accordance with manufacturers instructions. Service the Ranger pressure
infusor every six months or whenever service is needed. For replacement of the Ranger pressure infusor door latch, door, bladder, fingers, or power cord,
contact a biomedical technician. For additional technical support refer to the preventive maintenance manual or contact 3M.
Standby/ON mode
Condition
Cause
Solution
Nothing illuminates on the pressure
Power cord is not plugged into the power
Make sure the power cord is plugged into the power entry
infusor control panel when the main
entry module, or power cord is not plugged
module of the pressure infusor. Make sure the pressure
power switch is turned ON.
into an appropriately grounded outlet.
infusor is plugged into a properly grounded outlet.
Burned out LED light(s).
Contact a biomedical technician.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
Power LED status lights do not illuminate.
Power cord is not plugged into the power
Make sure the power cord is plugged into the power entry
entry module, or power cord is not plugged
module of the pressure infusor. Make sure the pressure
into an appropriately grounded outlet.
infusor is plugged into a properly grounded outlet.
Unit is not turned ON.
Using the main power switch located under the pressure
infusor turn the unit ON.
Burned out LED light.
Push the pressure infusor start/stop button, if unit functions
properly, continue use. Contact a biomedical technician
after use to replace LED.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
The status indicators (Low, In Range,
Power cord is not plugged into the power
Make sure the power cord is plugged into the power entry
and/or High) do not illuminate when the
entry module, or power cord is not plugged
module of the pressure infusor. Make sure the pressure
pressure infusor start/stop button is
into an appropriately grounded outlet.
infusor is plugged into a properly grounded outlet.
pushed.
Unit is not turned ON.
Using the main power switch located under the pressure
infusor turn the unit ON.
Burned out LED light(s).
Contact a biomedical technician.
Unit failure.
Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
Pressure infusor
Condition
Cause
Solution
Pressure infusor is not working.
Power cord is not plugged into the power entry
Make sure the power cord is plugged into the power entry
module, or power cord is not plugged into an
module of the pressure infusor. Make sure the pressure infusor is
appropriately grounded outlet.
plugged into a properly grounded outlet.
Unit is not turned ON.
Using the main power switch located under the pressure infusor
turn the unit ON.
Unit fault.
Discontinue use of unit. Contact a biomedical technician.
Blown fuse.
Contact a biomedical technician.
Low indicator
Pressure infusor bladder is loose or has become
Discontinue use of pressure infusor chamber. Use the other side
(solid yellow visual with audible
unattached.
of the pressure infusor.
indicator).
Reattach bladder by using your thumbs to fit one side of the bladder
port on the bladder retaining collar and stretch into position.
Pressure infusor door may not be closed and
Securely close and latch the pressure infusor door.
securely latched.
Detected pressure has fallen below 230 mmHg.
Continue infusion or use the other side of the pressure infusor.
Contact a biomedical technician after use.
High indicator (solid yellow
Pressure is above 330 mmHg.
Discontinue use of pressure infusor chamber. Use the other side of
visual with audible indicator).
the pressure infusor. Contact a biomedical technician after use.
Leakage of fluid.
Bag is not spiked securely.
Secure spike in bag.
Bladder does not deflate after
Unit fault.
Contact a biomedical technician after use.
pressure is discontinued.
7
Section 5: Maintenance and Storage
General Maintenance and Storage
All repairs and maintenance should be in accordance with manufacturers instructions. Service the Ranger pressure infusor every six months or whenever
service is needed. For replacement of the Ranger pressure infusor door latch, door, bladder, fingers, or power cord, contact a biomedical technician. For
additional technical support refer to the preventive maintenance manual or contact 3M.
Cleaning Instructions
1. Disconnect the Ranger unit from the power source before cleaning.
2. Cleaning should be performed in accordance with hospital practices for cleaning OR equipment. After every use; wipe the warming unit, the outside of
the warming unit hose, and any other surfaces that may have been touched. Use a damp, soft cloth and a hospital approved mild detergent, germicidal
disposable wipes, disinfecting towelettes, or antimicrobial spray. The following active ingredients are acceptable for use in cleaning the warming unit:
a. Oxidizers (e.g. 10% Bleach)
b. Quaternary Ammonium Compounds (e.g. 3M Quat Disinfectant Cleaner)
C. Phenolics (e.g. 3MTM Phenolic Disinfectant Cleaner)
d. Alcohols (e.g. 70% Isopropyl Alcohol
NOTICE:
1.
Do not use cleaning solutions with greater than 80% alcohol or solvents, including acetone and thinner, to clean the warming unit or hose. Solvents may
damage the labels and other plastic parts.
2. Do not immerse the Ranger unit or accessories in any liquid or subject them to any sterilization process.
Storage
Cover and store all components in a cool, dry place when not in use. Take care not to drop or jar the unit.
Servicing
All repairs and maintenance should be in accordance with manufacturers instructions. Service the Ranger pressure infusor every six months or whenever
service is needed. For replacement of the Ranger pressure infusor door latch, door, bladder, fingers, or power cord, contact a biomedical technician. For
additional technical support refer to the preventive maintenance manual or contact 3M.
Please report a serious incident occurring in relation to the device to 3M and the local competent authority (EU) or local regulatory authority.
Symbol Glossary
The following symbols may appear on the products labeling or exterior packaging.
"OFF" (power)
To indicate disconnection from the mains, at least for main switches, or their positions, and all those cases
where safety is involved. Source: IEC 60417-5008
"ON" (power)
To indicate connection to the mains, at least for mains switches, or their positions, and all those cases where
safety is involved. Source: IEC 60417-5007
Authorized Representative in
EC
REP
Indicates the authorized representative in the European Community. Source: ISO 15223, 5.1.2, 2014/35/EU,
European Community
and/or 2014/30/EU
Catalogue number
REF
Indicates the manufacturers catalogue number so that the medical device can be identified. Source: ISO 15223,
5.1.6
Caution
Indicates the need for the user to consult the instructions for use for important cautionary information such
as warnings and precautions that cannot, for a variety of reasons, be presented on the medical device itself.
Source: ISO 15223, 5.4.4
CE Mark 2797
C
€
Indicates conformity to all applicable European Union Regulations and Directives with notified body
2797
involvement.
Date of
Indicates the date when the medical device was manufactured. Source: ISO 15223, 5.1.3
Manufacture
Defibrillation-proof Type BF
Indicates the device applied part is Defibrillation-Proof Type BF. Source: IEC 60417-5334
applied part
Equipotentiality
To identify the terminals which, when connected together, bring the various parts of an equipment or of a
system to the same potential, not necessarily being the earth (ground) potential.
Source: IEC 60417-5021
Follow instructions for use
To signify that the instructions for use must be followed. Source: ISO 7010-M002
Fuse
Indicates a replaceable fuse
Importer
Indicates the entity importing the medical device into the EU
IP Code
IPX1
Indicates product can resist water that drips vertically onto it. Source: IEC 60529+AMD1:1999+
AMD2:2013CSV/COR2:2015
8
Keep dry
Indicates a medical device that needs to be protected from moisture. Source: ISO 15223, 5.3.4
Manufacturer
Indicates the medical device manufacturer as defined in EU Directives 90/385/EEC, 93/42/ EEC and 98/79/
EC. Source: ISO 15223, 5.1.1
Maximum safe working
Indicates the maximum safe working load in less than the number reported.
load
s25kg
Medical Device
MD
Indicates the item is a medical device
Protective earth; protective
To identify any terminal which is intended for connection to an external conductor for protection against electric
ground
shock in case of fault, or the terminal of a protective earth (ground) electrode. Source: IEC 60417-5019
Pushing prohibited
Indicates the device should not be pushed. Source: ISO 7010-P017
Recycle electronic
DO NOT throw this unit into a municipal trash bin when this unit has reached the end of its lifetime. Please
equipment
recycle. Source: Directive 2012/19/EC on waste electrical and electronic equipment (WEEE)
Rx Only
Indicates that U.S. Federal Law restricts this device to sale by or on the order of physician. 21 Code of Federal
Rx Only
Regulations (CFR) sec. 801.109(b)(1).
Serial number
SN
Indicates the manufacturers serial number so that a specific medical device can be identified. Source: ISO 15223, 5.1.7
UL Classified
CLASSIFIED
Indicates product was evaluated and Listed by UL for the USA and Canada.
US
Unique device identifier
UDI
Indicates bar code to scan product information into patient electronic health record
Section 6: Specifications
Guidance and manufacturers declaration - electromagnetic emissions
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the model 145 should assure that it is used
in such an environment.
Emissions test
Compliance
Electromagnetic environment - guidance
RF emissions
Group 1
The model 145 uses RF energy only for its internal function. Therefore, its RF emissions
CISPR 11
are very low and are not likely to cause any interference in nearby electronic equipment.
RF emissions
Class B
The model 145 is suitable for use in all establishments, including domestic
CISPR 11
establishments and those directly connected to the public low voltage power supply
network that supplies buildings used for domestic purposes.
Harmonic emissions
Class A
IEC 61000-3-2
Voltage fluctuations/Flicker emissions
Complies
IEC 61000-3-3
Guidance and manufacturers declaration - electromagnetic immunity
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the model 145 should assure that it is used
in such an environment.
Immunity Test
IEC 60601 test level
Compliance level
Electromagnetic environment - guidance
Electrostatic discharge (ESD)
+8 kV contact
+8 kV contact
Floors should be wood, concrete or ceramic tile. If floors are covered
IEC 61000-4-2
15 kV air
15 kV air
with synthetic material, the relative humidity should be at least 30%
Electrical fast transient/burst
2 kV power lines
2 kV power lines
Mains power quality should be that of a typical commercial or
IEC 61000-4-4
hospital environment.
Surge
1 kV line to line
1 kV line to line
Mains power quality should be that of a typical commercial or
IEC 61000-4-5
2 kV line to gnd
2 kV line to gnd
hospital environment.
Voltage dips, short interruptions
<5% U-
<5% U-
Mains power quality should be that of a typical commercial
and voltage variations on power
(>95% dip in UT)
(>95% dip in U-)
or hospital environment. If the user of the model 145 requires
supply input lines
for 0,5 cycle
for 0,5 cycle
continued operation during power mains interruptions, it
IEC 61000-4-11
40% U1
40% U1
is recommended that the model 145 be powered from an
(60% dip in UT)
(60% dip in U-)
uninterruptible power supply or a battery.
for 6 cycles
for 6 cycles
70% U1
70% U-
(30% dip in UT)
(30% dip in UT)
for 30 cycles
for 30 cycles
<5% U1
<5% U-
(>95% dip in UT)
(>95% dip in UT)
for 5 sec
for 5 sec
Power frequency (50/60 Hz)
3 A/m
3 A/m
Power frequency magnetic fields should be at levels characteristic
magnetic field
of a typical location in a typical commercial or hospital environment.
IEC 61000-4-8
NOTE UT is the a.c. mains voltage prior to application of the test level.
9
Guidance and manufacturers declaration - electromagnetic immunity
The model 145 is intended for use in the electromagnetic environment specified below. The customer or the user of the model 145 should assure that it is used
in such an environment.
Immunity Test
IEC 60601 test level
Compliance level
Electromagnetic environment - guidance
Portable and mobile RF communications equipment should be used no closer to any
part of the model 145, including cables, than the recommended separation distance
calculated from the equation applicable to the frequency of the transmitter.
Recommended separation distance
d = 1.2 VP
Conducted RF
3 Vrms
d = 1.2 VP 80 MHz to 800 MHz
IEC 61000-4-6
150 kHz to 80 MHz
3 Vrms
d = 2.3 VP 800 MHz to 2,5 GHz
where P is the maximum output power rating of the transmitter in watts (W) according to
Radiated RF
3 V/m
3 V/m
the transmitter manufacturer and d is the recommended separation distance in meters (m).
IEC 61000-4-3
80 MHz to 2,5 GHz
Field strengths from fixed RF transmitters, as determined by an electromagnetic site
survey, ,should be less than the compliance level in each frequency rangeb.
Interference may occur in the vicinity of equipment marked with the following symbol:
(((.)))
NOTE 1 At 80 MHz and 800 MHz, the higher frequency range applies.
NOTE 2 These guidelines may not apply in all situations. Electromagnetic propagation is affected by absorption and reflection from structures, objects and people.
Field strengths from fixed transmitters, such as base stations for radio (cellular/cordless) telephones and land mobile radios, amateur radio, AM and
FM radio broadcast and TV broadcast cannot be predicted theoretically with accuracy. To assess the electromagnetic environment due to fixed RF
transmitters, an electromagnetic site survey should be considered. If the measured field strength in the location in which the model 145 is used exceeds
the applicable RF compliance level above, the model 145 should be observed to verify normal operation. If abnormal performance is observed, additional
measures may be necessary, such as re-orienting or relocating the model 145.
b
Over the frequency range 150 kHz to 80 MHz, field strengths should be less than 3 V/m.
Recommended separation distances between portable and mobile RF communications equipment and the model 145
The model 145 is intended for use in an electromagnetic environment in which radiated RF disturbances are controlled. The customer or the user of the
model 145 can help prevent electromagnetic interference by maintaining a minimum distance between portable and mobile RF communications equipment
(transmitters) and the model 145 as recommended below, according to the maximum output power of the communications equipment.
Separation distance according to frequency of transmitter
Rated maximum output
m
power of transmitter
W
150 kHz to 80 MHz
80 MHz to 800 MHz
800 MHz to 2,5 GHz
d = 1.2 VP
d = 1.2 VP
d = 2.3 VP
0,01
0.12
0.12
0.23
0,1
0.37
0.37
0.74
1
1.17
1.17
2.30
10
3.69
3.69
7.37
100
11.67
11.67
23.30
For transmitters rated at a maximum output power not listed above, the recommended separation distance d in meters (m) can be estimated using the
equation applicable to the frequency of the transmitter, where P is the maximum output power rating of the transmitter in watts (W) according to the
transmitter manufacturer.
NOTE 1 At 80 MHz and 800 MHz, the separation distance for the higher frequency range applies.
NOTE 2 These guidelines may not apply in all situations. Electromagnetic propagation is affected by absorption and reflection from structures, objects and
people.
Physical characteristics
Dimensions
15.75 in (40 cm) high
20 in (51 cm) wide, 7.75 in (20 cm) deep
Weight
17 lb. (7.7 kg)
Mounting
Dual clamp
10
Classification
Protection against electric shock:
Class I Medical Electrical Equipment with Type BF defibrillation-proof applied parts
Protection against ingress of water: IPX1
Mode of operation: Continuous operation.
CLASSIFIED
MEDICAL - GENERAL MEDICAL EQUIPMENT AS TO ELECTRIC SHOCK, FIRE AND MECHANICAL HAZARDS ONLY IN ACCORDANCE
C
UL
WITH ANSI/AAMI ES 60601-1 (2005) + AMD (2012), CAN/CSA-C22.2 No. 60601-1 (2008) + (2014), and IEC 60601-1-6:2010 (Third
US
Edition) + A1:2013; Control No. 4HZ8
Electrical characteristics
Leakage current
Meets leakage current requirements in accordance with IEC 606011.
Power cord
15 feet (4.6 m)
Device rating
110-120 VAC, 50/60 Hz, 1 Amp
220-240 VAC, 50/60 Hz, 0.8 Amp
Fuse
2 X F1A-H, rated 250V, for 110-120 VAC unit
2 X F0.8A-H, rated 250V, for 220-240 VAC unit
Storage and transport conditions
Storage/transport temperature
-20 to 60°C (-4°F to 140°F)
Store all components at room temperature and in a dry place when in use.
Operating humidity
Up to 90% RH, noncondensing
Atmospheric pressure range
Altitude up to 2000m or 80 kPa
Performance characteristics
Operating pressure
300 + 10 mmHg setpoint
Note:
Pressure system is In Range when the pressure infusor bladders are inflated to between 230 mmHg (low) and 330 mmHg (high). If pressure falls below
230 mmHg for more than approximately 30 seconds the Low yellow indicator will illuminate and an audible indicator will sound. The High yellow and
audible indicator notifies the user when the pressure infusor bladder is above 330 mmHg.
The outlet pressure of the fluid may vary with the surface area and volume of the fluid bag.
3M and Ranger are trademarks of 3M.
EC REP
3M Deutschland GmbH
MD
2797
Used under license in Canada.
Health Care Business
Made in USA with Globally
© 2020, 3M. Unauthorized use prohibited.
Carl-Schurz-Str. 1
Sourced Materials
All rights reserved.
41453 Neuss, Germany
3M Company
3M et Ranger sont des marques de commerce de 3M.
2510 Conway Ave.
Utilisées sous licence au Canada.
St. Paul, MN 55144 USA
© 2020, 3M. Toute utilisation non autorisée est interdite.
Issue Date 2020-09
1-800-228-3957 (USA Only)
Tous droits réservés.
34-8726-1613-0
</document2>

History: {{ .History }}
input: {{ .Input }}

Please confirm your understanding of the above requirements

Assistant: I understand the requirements.
I will read the english manuals.
Differences that I should point out are related to information about use, cleaning, care, and maintenance of the equipment mentioned in the manual.
I will inform of any differences I find, and I will explain what changed and where.

Human: Readback is correct, please tell me what differences you see.
Assistant:
`

func buildPrompt() string {
	//prompt := fmt.Sprintf(DEFAULT_PROMPT_TEMPLATE, "no history", "What is the best way to prompt an llm?")
	template, err := template.New("prompt").Parse(DEFAULT_PROMPT_TEMPLATE)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := struct {
		History string
		Input   string
	}{
		History: "no history",
		Input:   "What is the best way to prompt an llm?",
	}

	var builder strings.Builder

	err = template.Execute(&builder, data)

	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return builder.String()
}

func runRequest() {
	ctx := context.Background()
	llm, err := bedrock.New(
		bedrock.WithModel("anthropic.claude-v2:1"),
	)

	if err != nil {
		log.Fatalf("Failed to create LLM: %v", err)
	}

	prompt := buildPrompt()

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt,
		llms.WithMaxTokens(30000),
		llms.WithMaxLength(20000),
	)

	if err != nil {
		log.Fatalf("Failed to generate completion: %v", err)
	}

	log.Println(completion)
}

func handler() {
	runRequest()
}

func main() {
	log.Println("Starting LLM...")
	lambda.Start(handler)
	log.Println("LLM completed.")
}
