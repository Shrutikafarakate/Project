import os
import numpy as np
from PIL import Image
from django.shortcuts import render, redirect
from django.core.mail import send_mail
from django.conf import settings
from django.core.files.storage import FileSystemStorage
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.models import User
from django.contrib import messages
from django.contrib.auth.decorators import login_required

from .forms import RecommendationForm, ContactForm, UploadImageForm
from .models import UserProfile
from .utils import generate_pdf_report  # PDF report utility

import tensorflow as tf

# ------------------ USER AUTH ------------------

def user_signup(request):
    if request.method == "POST":
        username = request.POST.get("username")
        email = request.POST.get("email")
        password = request.POST.get("password")
        confirm_password = request.POST.get("confirm_password")
        dob = request.POST.get("dob")
        location = request.POST.get("location")
        phone = request.POST.get("phone")

        if password == confirm_password:
            if User.objects.filter(username=username).exists():
                messages.error(request, "Username already exists!")
            elif User.objects.filter(email=email).exists():
                messages.error(request, "Email is already in use!")
            else:
                user = User.objects.create_user(username=username, email=email, password=password)
                user.save()
                request.session["dob"] = dob
                request.session["location"] = location
                messages.success(request, "Account created successfully! Please log in.")
                return redirect("login")
        else:
            messages.error(request, "Passwords do not match!")

    return render(request, "signup.html")


def user_login(request):
    if request.method == "POST":
        username = request.POST.get("username")
        password = request.POST.get("password")
        user = authenticate(request, username=username, password=password)

        if user is not None:
            login(request, user)
            return redirect("index")
        else:
            messages.error(request, "Invalid username or password!")

    return render(request, "login.html")


def user_logout(request):
    logout(request)
    return redirect("login")


# ------------------ MODEL LOADING ------------------

# 🔧 Change this path to your local model location
MODEL_PATH = r"C:\Users\shrut\Downloads\Archives\Skin-Lesion-Detection-main (1)\Skin-Lesion-Detection-main\detection\models\HAM10000_CNN.h5"

model = None
if os.path.exists(MODEL_PATH):
    try:
        model = tf.keras.models.load_model(MODEL_PATH)
        print(f"✅ Model loaded successfully from {MODEL_PATH}")
    except Exception as e:
        print(f"⚠️ Error loading model: {e}")
        model = None
else:
    print(f"⚠️ Model file NOT found at {MODEL_PATH}. Predictions disabled.")


# ------------------ CLASS LABELS ------------------

# ⚠️ Ensure this matches your training order
lesion_classes = {
    0: "Actinic Keratoses (AKIEC)",
    1: "Basal Cell Carcinoma (BCC)",
    2: "Benign Keratosis (BKL)",
    3: "Dermatofibroma (DF)",
    4: "Melanoma (MEL)",
    5: "Nevus (NV)",
    6: "Vascular Lesions (VASC)"
}


# ------------------ STATIC PAGES ------------------

def index(request):
    return render(request, 'index.html')


def about(request):
    return render(request, 'about.html')


def faq(request):
    return render(request, 'faq.html')


# ------------------ IMAGE UPLOAD & PREDICTION ------------------

@login_required
def home(request):
    form = UploadImageForm()
    prediction = None
    uploaded_file_url = None
    report_url = None

    if request.method == 'POST' and request.FILES.get('image'):
        if model is None:
            prediction = "⚠️ Model not available. Please contact the administrator."
        else:
            image = request.FILES['image']
            fs = FileSystemStorage()
            filename = fs.save(image.name, image)
            uploaded_file_url = fs.url(filename)
            img_path = fs.path(filename)

            try:
                # ✅ Image Preprocessing (match with training)
                img = Image.open(img_path).convert("RGB")
                img = img.resize((224, 224))  # use same size as training
                img = np.array(img) / 255.0
                img = np.expand_dims(img, axis=0)

                predictions = model.predict(img)
                predicted_class = np.argmax(predictions)
                confidence = float(np.max(predictions)) * 100
                
                # 🎯 Confidence threshold check
                CONFIDENCE_THRESHOLD = 70.0
                
                if confidence < CONFIDENCE_THRESHOLD:
                    prediction = "❌ This does not appear to be a recognizable skin lesion. Please upload a clear image of a skin lesion."
                else:
                    prediction = lesion_classes.get(predicted_class, "Unknown")
                    print(f"🧠 Raw Predictions: {predictions}")
                    print(f"✅ Predicted Class: {predicted_class} -> {prediction} ({confidence:.2f}% confidence)")

                    # Generate PDF report without DOB and location
                    report_url = generate_pdf_report(request.user, img_path, prediction)

            except Exception as e:
                prediction = f"⚠️ Error processing image: {str(e)}"

    return render(request, 'home.html', {
        'form': form,
        'uploaded_file_url': uploaded_file_url,
        'prediction': prediction,
        'report_url': report_url
    })

# ------------------ CONTACT ------------------

def contact(request):
    if request.method == "POST":
        form = ContactForm(request.POST)
        if form.is_valid():
            name = form.cleaned_data["name"]
            email = form.cleaned_data["email"]
            message = form.cleaned_data["message"]

            subject = f"New Contact Form Submission from {name}"
            body = f"Name: {name}\nEmail: {email}\n\nMessage:\n{message}"

            send_mail(
                subject,
                body,
                settings.DEFAULT_FROM_EMAIL,
                ['shrutikafarakate492@gmail.com'],
                fail_silently=False,
            )
            return render(request, "thank_you.html")
    else:
        form = ContactForm()

    return render(request, "contact.html", {"form": form})


# ------------------ RECOMMENDATIONS ------------------

recommendations = {
    'mel': (
        "For melanoma (MEL):\n"
        "- Consult a dermatologist immediately.\n"
        "- Regularly check your skin for new or changing moles.\n"
        "- Avoid sun exposure and wear sunscreen with SPF 30+."
    ),
    'bcc': (
        "For basal cell carcinoma (BCC):\n"
        "- Early detection and removal by a healthcare professional are key.\n"
        "- Avoid sun exposure and use sunscreen daily."
    ),
    'akiec': (
        "For actinic keratosis (AKIEC):\n"
        "- Consult a dermatologist for treatment options like cryotherapy.\n"
        "- Use sunscreen regularly."
    ),
    'df': (
        "For dermatofibroma (DF):\n"
        "- Typically harmless, but if you notice changes, consult a dermatologist."
    ),
    'nv': (
        "For benign nevi (NV):\n"
        "- Regularly monitor any moles for changes in size, shape, or color."
    ),
    'bkl': (
        "For benign keratosis (BKL):\n"
        "- Usually harmless but should be checked by a dermatologist."
    ),
    'vasc': (
        "For vascular lesions (VASC):\n"
        "- Consult a healthcare provider if the lesion is growing or bleeding."
    ),
}


def recommendation(request):
    recommendation_message = ""

    if request.method == 'POST':
        form = RecommendationForm(request.POST)
        if form.is_valid():
            skin_condition = form.cleaned_data['skin_condition']
            recommendation_message = recommendations.get(
                skin_condition, "No recommendation available for this condition."
            )
    else:
        form = RecommendationForm()

    return render(request, 'recommendation.html', {
        'form': form,
        'recommendation_message': recommendation_message
    })